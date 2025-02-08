package model

func NewLoadingModel() *LoadingModel {
	return &LoadingModel{}
}

type LoadingModel struct {
	state LoadingState
}

func (l *LoadingModel) State() LoadingState {
	return l.state
}

func (l *LoadingModel) SetState(state LoadingState) {
	l.state = state
}

type LoadingState struct {
	Promise         LoadingPromise
	SuccessViewName ViewName
	ErrorViewName   ViewName
}
