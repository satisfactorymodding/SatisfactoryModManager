package bindings

import (
	"context"

	"github.com/satisfactorymodding/SatisfactoryModManager/project_file"
)

type ProjectFile struct {
	ctx context.Context
}

func MakeProjectFile() (*ProjectFile, error) {
	pf := &ProjectFile{}

	return pf, nil
}

func (s *ProjectFile) startup(ctx context.Context) {
	s.ctx = ctx
}

func (s *ProjectFile) GetProjectFile() project_file.Project {
	return project_file.ProjectFile
}
