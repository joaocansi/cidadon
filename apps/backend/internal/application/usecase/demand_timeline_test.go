package usecase

import (
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"
	"testing"
)

func TestValidateTimelineInputRequiresWrittenJustification(t *testing.T) {
	tests := []struct {
		name  string
		input service.DemandTimelineInput
		valid bool
	}{
		{name: "empty", input: service.DemandTimelineInput{}, valid: false},
		{name: "two characters", input: service.DemandTimelineInput{Message: "ok"}, valid: false},
		{name: "message", input: service.DemandTimelineInput{Message: "Equipe técnica acionada."}, valid: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateTimelineInput(test.input)
			if (err == nil) != test.valid {
				t.Fatalf("validateTimelineInput() error = %v, valid = %v", err, test.valid)
			}
		})
	}
}

func TestPublicTimelineExcludesCommentAndModerationAuditEvents(t *testing.T) {
	visible := map[string]bool{}
	for _, eventType := range publicTimelineEventTypes() {
		visible[eventType] = true
	}
	for _, eventType := range []string{"commented", "comment_deleted", "comment_hidden", "comment_reported"} {
		if visible[eventType] {
			t.Fatalf("%q must remain outside the public timeline", eventType)
		}
	}
	if !visible["milestone"] || !visible["automatically_completed"] {
		t.Fatal("operational milestones and automatic completion must be visible")
	}
}

func TestStatusMetadataPreservesMessageImagesAndTransition(t *testing.T) {
	input := service.DemandTimelineInput{Message: "A equipe iniciou a vistoria.", Images: []string{"data:image/png;base64,AA=="}}
	metadata := statusMetadata(&input, entity.DemandStatusReview, entity.DemandStatusInProgress)
	if metadata["message"] != input.Message {
		t.Fatalf("message = %v", metadata["message"])
	}
	if metadata["from_status"] != entity.DemandStatusReview || metadata["to_status"] != entity.DemandStatusInProgress {
		t.Fatalf("unexpected transition metadata: %#v", metadata)
	}
	images, ok := metadata["images"].([]string)
	if !ok || len(images) != 1 {
		t.Fatalf("images = %#v", metadata["images"])
	}
}
