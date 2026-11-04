// An entry that has been delivered as many times as the group tolerates is
// parked instead of handled again.
export function shouldDeadLetter(
	deliveries: number,
	maxDeliveries: number,
): boolean {
	return deliveries > maxDeliveries
}
