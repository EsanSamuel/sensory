package helpers

import (
	"encoding/json"
	"os"
)

type ProjectMeta struct {
	ProjectName string `json:"project_name"`
	ProjectId   string `json:"project_id"`
	UserId      string `json:"user_id"`
}

func SaveProjectMeta(projectName, projectId, userId string) error {
	meta := ProjectMeta{
		ProjectName: projectName,
		ProjectId:   projectId,
		UserId:      userId,
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(".sensory.json", data, 0644)
}
