package configuration

// Version is the fence version this contract reads. Any other version was
// written for another shape and is refused rather than misread.
const Version = 2

// replan is what every refusal of another fence shape asks for.
const replan = "re-plan the runtime"
