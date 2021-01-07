package lostfiles

import (
	"solt/msvc"
)

type projecter func(*msvc.MsbuildProject)

func enumerate(projects []*msvc.MsbuildProject, handlers ...projecter) {
	for _, prj := range projects {
		for _, handler := range handlers {
			handler(prj)
		}
	}
}
