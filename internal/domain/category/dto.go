package category

import (
	"errors"
	"net/http"
)

type Request struct {
	ParentID string `json:"parentID,omitempty"`
	Name     string `json:"name"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Name == "" {
		return errors.New("name: cannot be blank")
	}

	return nil
}

type Response struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Children []Response `json:"children,omitempty"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:   data.ID,
		Name: *data.Name,
	}

	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
