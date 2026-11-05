package configuration

import (
	"slices"
	"strings"
)

// componentProblems are the problems of one component on its own: its path,
// stack, role, port and image, its share and its replica count.
func componentProblems(index int, component FenceComponent) []Problem {
	insideLangdrift := component.Path == ".langdrift" || strings.HasPrefix(component.Path, ".langdrift/")

	var problems []Problem

	if CheckRelativePath(component.Path) != nil {
		problems = append(problems, problemAt(componentField(index, "path"), "must be a clean relative path"))
	} else if insideLangdrift {
		problems = append(problems, problemAt(componentField(index, "path"), "must be outside .langdrift"))
	}

	if !slices.Contains(FenceStacks, component.Stack) {
		problems = append(problems, problemAt(componentField(index, "stack"), "must be one of "+strings.Join(FenceStacks, ", ")))
	}

	if !slices.Contains(FenceRoles, component.Role) {
		problems = append(problems, problemAt(componentField(index, "role"), "must be one of "+strings.Join(FenceRoles, ", ")))
	}

	if !ValidPort(component.Port) {
		problems = append(problems, problemAt(componentField(index, "port"), "must be 1 to 65535"))
	}

	if component.Image == "" || !IsImageReference(component.Image) {
		problems = append(problems, problemAt(componentField(index, "image"), "must be an image reference"))
	}

	if component.CPUMillis < 1 {
		problems = append(problems, problemAt(componentField(index, "cpuMillis"), "must be a positive share"))
	}

	if component.MemoryMiB < 1 {
		problems = append(problems, problemAt(componentField(index, "memoryMiB"), "must be a positive share"))
	}

	if component.Replicas != 1 {
		problems = append(problems, problemAt(componentField(index, "replicas"), "must be 1"))
	}

	return problems
}
