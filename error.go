package pixiv

import "fmt"

type ErrToken struct {
	StatusCode int
	Status     string
	Body       TokenErrorBody
}

func (e ErrToken) Error() string {
	return fmt.Sprintf("%s", e.Status)
}

type ErrAPI struct {
	StatusCode int
	Status     string
	Body       APIErrorBody
}

func (e ErrAPI) Error() string {
	return fmt.Sprintf("%s", e.Status)
}

type ErrInvalidParams struct {
	Errs []ErrInvalidParam
}

func (e *ErrInvalidParams) Error() string {
	return fmt.Sprintf("%d validation error(s) found", len(e.Errs))
}

func (e *ErrInvalidParams) Add(err ErrInvalidParam) {
	e.Errs = append(e.Errs, err)
}

func (e *ErrInvalidParams) Len() int {
	return len(e.Errs)
}

type ErrInvalidParam struct {
	Field   string
	Message string
}

func (e ErrInvalidParam) Error() string {
	return fmt.Sprintf("%s, %s", e.Field, e.Message)
}
