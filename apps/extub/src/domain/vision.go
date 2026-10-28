package domain

type OverviewSummary struct {
	FixtureVersion       string `json:"fixtureVersion"`
	Illustrative         bool   `json:"illustrative"`
	ProductVision        uint32 `json:"productVision"`
	IntentionalEvolution uint32 `json:"intentionalEvolution"`
	UnexplainedDrift     uint32 `json:"unexplainedDrift"`
}

type FixtureDate struct {
	Month uint32 `json:"month"`
	Day   uint32 `json:"day"`
}

type ActorFixture struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Initials  string `json:"initials"`
	AvatarKey string `json:"avatarKey"`
	TeamKey   string `json:"teamKey"`
}

type EventFixture struct {
	ID             string         `json:"id"`
	Actors         []ActorFixture `json:"actors"`
	Classification string         `json:"classification"`
	Date           FixtureDate    `json:"date"`
	DateKey        string         `json:"dateKey"`
	DecisionKey    string         `json:"decisionKey"`
	Delta          int32          `json:"delta"`
	ProductAreaKey string         `json:"productAreaKey"`
	ReasonKey      string         `json:"reasonKey"`
	TitleKey       string         `json:"titleKey"`
}

type HistoryPoint struct {
	Month   string        `json:"month"`
	Value   uint32        `json:"value"`
	EventID string        `json:"eventId"`
	Event   *EventFixture `json:"event,omitempty"`
}

type OverviewHistory struct {
	FixtureVersion string         `json:"fixtureVersion"`
	Illustrative   bool           `json:"illustrative"`
	Points         []HistoryPoint `json:"points"`
}

type VisionFacts struct {
	BaselineValue        uint32 `json:"baselineValue"`
	CurrentValue         uint32 `json:"currentValue"`
	BaselineDelta        uint32 `json:"baselineDelta"`
	BaselineVersion      string `json:"baselineVersion"`
	MovementCount        uint32 `json:"movementCount"`
	IntentionalEvolution uint32 `json:"intentionalEvolution"`
	UnexplainedDrift     uint32 `json:"unexplainedDrift"`
}

type BaselineFixture struct {
	Version         string       `json:"version"`
	Date            FixtureDate  `json:"date"`
	DateKey         string       `json:"dateKey"`
	Approver        ActorFixture `json:"approver"`
	ApproverRole    string       `json:"approverRole"`
	SourceArtifact  string       `json:"sourceArtifact"`
	StrategyVersion uint32       `json:"strategyVersion"`
	ProductID       string       `json:"productId"`
	ScopeKey        string       `json:"scopeKey"`
	AreaKeys        []string     `json:"areaKeys"`
	SupersedesKey   string       `json:"supersedesKey"`
	ReasonKey       string       `json:"reasonKey"`
}

type DecisionParticipant struct {
	Actor     ActorFixture `json:"actor"`
	ActionKey string       `json:"actionKey"`
}

type EvidenceReference struct {
	Artifact   string `json:"artifact"`
	MessageKey string `json:"messageKey"`
	ActorName  string `json:"actorName"`
}

type DecisionFixture struct {
	ID           string                `json:"id"`
	TitleKey     string                `json:"titleKey"`
	Date         FixtureDate           `json:"date"`
	DateKey      string                `json:"dateKey"`
	AreaKey      string                `json:"areaKey"`
	Delta        int32                 `json:"delta"`
	PeopleKey    string                `json:"peopleKey"`
	Participants []DecisionParticipant `json:"participants"`
	Evidence     []EvidenceReference   `json:"evidence"`
}

type PersonFixture struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	FullName        string `json:"fullName"`
	Initials        string `json:"initials"`
	AvatarKey       string `json:"avatarKey"`
	TeamKey         string `json:"teamKey"`
	RoleKey         string `json:"roleKey"`
	RoleLiteral     string `json:"roleLiteral"`
	DetailKey       string `json:"detailKey"`
	Decisions       uint32 `json:"decisions"`
	Approvals       uint32 `json:"approvals"`
	Implementations uint32 `json:"implementations"`
	Reviews         uint32 `json:"reviews"`
	Proposals       uint32 `json:"proposals"`
	OwnsKey         string `json:"ownsKey"`
}

type TeamFixture struct {
	Key         string   `json:"key"`
	AreaKeys    []string `json:"areaKeys"`
	AreasKey    string   `json:"areasKey"`
	DetailKey   string   `json:"detailKey"`
	Decisions   uint32   `json:"decisions"`
	Unexplained uint32   `json:"unexplained"`
	Approvals   uint32   `json:"approvals"`
}

type EvidenceFixture struct {
	Artifact       string `json:"artifact"`
	ObservationKey string `json:"observationKey"`
	StatusKey      string `json:"statusKey"`
}

type ReportFixture struct {
	FromValue            uint32 `json:"fromValue"`
	ToValue              uint32 `json:"toValue"`
	IntentionalEvolution uint32 `json:"intentionalEvolution"`
	UnexplainedDrift     uint32 `json:"unexplainedDrift"`
	SummaryKey           string `json:"summaryKey"`
}

type AttentionFixture struct {
	ID             string `json:"id"`
	EventID        string `json:"eventId"`
	Classification string `json:"classification"`
	TitleKey       string `json:"titleKey"`
	DetailKey      string `json:"detailKey"`
	ActorName      string `json:"actorName"`
	TeamKey        string `json:"teamKey"`
	ActionKey      string `json:"actionKey"`
}

type CommentFixture struct {
	Author     ActorFixture `json:"author"`
	Date       FixtureDate  `json:"date"`
	DateKey    string       `json:"dateKey"`
	MessageKey string       `json:"messageKey"`
}

type VisionProductFixture struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
}

type AccountFixture struct {
	ID        string `json:"id"`
	AvatarKey string `json:"avatarKey"`
}

type NotificationFixture struct {
	Count          uint32   `json:"count"`
	EventIDs       []string `json:"eventIds"`
	DescriptionKey string   `json:"descriptionKey"`
}

type VoiceFixture struct {
	BaselineValue     uint32   `json:"baselineValue"`
	CurrentValue      uint32   `json:"currentValue"`
	WeeklyIntentional uint32   `json:"weeklyIntentional"`
	WeeklyUnexplained uint32   `json:"weeklyUnexplained"`
	AnswerKeys        []string `json:"answerKeys"`
	FallbackKey       string   `json:"fallbackKey"`
}

type WorkspaceArtifact struct {
	ID         string   `json:"id"`
	Kind       string   `json:"kind"`
	TitleKey   string   `json:"titleKey"`
	SummaryKey string   `json:"summaryKey"`
	Path       string   `json:"path"`
	Source     string   `json:"source"`
	Status     string   `json:"status"`
	OccurredAt string   `json:"occurredAt"`
	Revision   string   `json:"revision"`
	References []string `json:"references"`
	ContentKey string   `json:"contentKey"`
}

type WorkspaceDrift struct {
	ID               string   `json:"id"`
	TitleKey         string   `json:"titleKey"`
	SummaryKey       string   `json:"summaryKey"`
	Classification   string   `json:"classification"`
	ReviewStatus     string   `json:"reviewStatus"`
	ProductionStatus string   `json:"productionStatus"`
	TargetID         string   `json:"targetId"`
	ArtifactIDs      []string `json:"artifactIds"`
	ObservedAt       string   `json:"observedAt"`
	DecisionID       string   `json:"decisionId"`
}

type WorkspaceTarget struct {
	ID          string   `json:"id"`
	TitleKey    string   `json:"titleKey"`
	SummaryKey  string   `json:"summaryKey"`
	Revision    string   `json:"revision"`
	Status      string   `json:"status"`
	ArtifactIDs []string `json:"artifactIds"`
}

type WorkspaceTimelineEvent struct {
	ID         string `json:"id"`
	TitleKey   string `json:"titleKey"`
	Date       string `json:"date"`
	Kind       string `json:"kind"`
	ArtifactID string `json:"artifactId"`
	DriftID    string `json:"driftId"`
}

type WorkspaceDiffFile struct {
	Path      string `json:"path"`
	Before    string `json:"before"`
	After     string `json:"after"`
	Additions uint32 `json:"additions"`
	Deletions uint32 `json:"deletions"`
}

type WorkspacePullRequest struct {
	ID         string              `json:"id"`
	Number     uint32              `json:"number"`
	Repository string              `json:"repository"`
	BaseSHA    string              `json:"baseSha"`
	HeadSHA    string              `json:"headSha"`
	Branch     string              `json:"branch"`
	Files      []WorkspaceDiffFile `json:"files"`
}

type WorkspaceFixture struct {
	Version     string                   `json:"version"`
	Artifacts   []WorkspaceArtifact      `json:"artifacts"`
	Drifts      []WorkspaceDrift         `json:"drifts"`
	Targets     []WorkspaceTarget        `json:"targets"`
	Timeline    []WorkspaceTimelineEvent `json:"timeline"`
	PullRequest *WorkspacePullRequest    `json:"pullRequest,omitempty"`
}

type VisionFixture struct {
	FixtureVersion string                 `json:"fixtureVersion"`
	Illustrative   bool                   `json:"illustrative"`
	Overview       OverviewSummary        `json:"overview"`
	History        OverviewHistory        `json:"history"`
	Vision         VisionFacts            `json:"vision"`
	Baseline       BaselineFixture        `json:"baseline"`
	Decisions      []DecisionFixture      `json:"decisions"`
	People         []PersonFixture        `json:"people"`
	Teams          []TeamFixture          `json:"teams"`
	Evidence       []EvidenceFixture      `json:"evidence"`
	Report         ReportFixture          `json:"report"`
	Attention      []AttentionFixture     `json:"attention"`
	Comment        CommentFixture         `json:"comment"`
	Products       []VisionProductFixture `json:"products"`
	Account        AccountFixture         `json:"account"`
	Notifications  NotificationFixture    `json:"notifications"`
	Voice          VoiceFixture           `json:"voice"`
	Workspace      WorkspaceFixture       `json:"workspace"`
}
