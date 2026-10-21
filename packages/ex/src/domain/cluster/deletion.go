package cluster

import "fmt"

// AuthorizeDeletion allows deleting the cluster only when everything in it is
// either part of the creation baseline or controlled, through a chain of
// controller references, by an unchanged baseline object. Foreign namespaces,
// replaced namespaces, vanished APIs, and unrecognized objects all block it.
// Evidence without a UID never identifies anything, so it blocks deletion too.
func (record *Record) AuthorizeDeletion(namespaces map[string]string, observed Inventory) error {
	if err := record.Baseline.matchNamespaces(namespaces); err != nil {
		return err
	}

	bootstrap := record.Baseline.Inventory
	if record.Baseline.InventoryVersion != 1 || len(bootstrap.Objects) == 0 || len(bootstrap.ResourceTypes) == 0 {
		return unsafeDeletion("complete bootstrap inventory is missing")
	}

	if err := requireResourceTypes(bootstrap.ResourceTypes, observed.ResourceTypes); err != nil {
		return err
	}

	known := bootstrap.recognize(observed)

	for key, identity := range observed.Objects {
		if identity.UID == "" || !known[identity.UID] {
			return unsafeDeletion("resource %q is foreign or changed", key)
		}
	}

	return nil
}

// matchNamespaces requires the observed namespaces to be exactly the baseline.
// A missing UID would compare equal to an absent map entry, so it never matches.
func (baseline Baseline) matchNamespaces(observed map[string]string) error {
	if len(baseline.Namespaces) == 0 {
		return unsafeDeletion("namespace baseline is missing")
	}

	for name, uid := range observed {
		if uid == "" || baseline.Namespaces[name] != uid {
			return unsafeDeletion("namespace %q is foreign or was replaced", name)
		}
	}

	for name, uid := range baseline.Namespaces {
		if uid == "" || observed[name] != uid {
			return unsafeDeletion("baseline namespace %q is missing, changed, or unidentified", name)
		}
	}

	return nil
}

func requireResourceTypes(recorded, observed []string) error {
	available := make(map[string]bool, len(observed))
	for _, resourceType := range observed {
		available[resourceType] = true
	}

	for _, resourceType := range recorded {
		if !available[resourceType] {
			return unsafeDeletion("resource API %q disappeared; inventory is incomplete", resourceType)
		}
	}

	return nil
}

// recognize returns the UIDs of observed objects that are unchanged bootstrap
// objects, or that a recognized object controls. The empty UID is never known.
func (bootstrap Inventory) recognize(observed Inventory) map[string]bool {
	known := map[string]bool{}

	for key, identity := range observed.Objects {
		if original, exists := bootstrap.Objects[key]; exists && identity.UID != "" && original.UID == identity.UID {
			known[identity.UID] = true
		}
	}

	for changed := true; changed; {
		changed = false

		for _, identity := range observed.Objects {
			if identity.UID == "" || known[identity.UID] {
				continue
			}

			for _, controller := range identity.Controllers {
				if known[controller] {
					known[identity.UID] = true
					changed = true

					break
				}
			}
		}
	}

	return known
}

func unsafeDeletion(format string, arguments ...any) error {
	return fmt.Errorf("%w: %s", ErrUnsafeDeletion, fmt.Sprintf(format, arguments...))
}
