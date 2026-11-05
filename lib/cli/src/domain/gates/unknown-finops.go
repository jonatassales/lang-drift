package gates

// noFinOpsSource is why CI reads no FinOps value: the inputs and the ledger
// live in .langdrift/, which stays on the machine that ran dft.
const noFinOpsSource = "no committed FinOps source: `.langdrift/` stays on the machine that ran `dft`"

// UnknownFinOps is the FinOps section of a report made without a FinOps
// source: the budget, the spend and the ledger, each unknown with the reason.
func UnknownFinOps() []Reading {
	var readings []Reading

	for _, name := range []string{"Budget", "Spend", "Ledger"} {
		readings = append(readings, Reading{
			Name:   name,
			Status: StatusUnknown,
			Reason: noFinOpsSource,
		})
	}

	return readings
}
