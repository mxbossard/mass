package compare

import (
	k8sv1 "k8s.io/api/core/v1"
)

type differ[T any] interface {
	Kind() string
	Old() T
	New() T
}

type diff[T any] struct {
	name string
	old  T
	new  T
}

func (d diff[T]) Kind() string {
	return d.name
}

func (d diff[T]) Old() T {
	return d.old
}

func (d diff[T]) New() T {
	return d.new
}

type podDiff struct {
	diff[k8sv1.Pod]

	initContainers []containerDiff
	containers     []containerDiff
	volumes        []volumeDiff
}

func (d podDiff) DoDiffer() bool {

}

func (d podDiff) IsRestartNeeded() bool {

}

func (d podDiff) UpdatableDiffs() []differ[any] {

}

func (d podDiff) Diffs() (diffs []differ[any]) {
	//var d1 differ[any]
	d1 := d.initContainers[0]

	diffs = append(diffs, d1.image)
}

type volumeDiff struct {
	diff[k8sv1.Volume]

	data diff[string]
}

type containerDiff struct {
	diff[any]

	image       diff[string]
	command     diff[[]string]
	args        diff[[]string]
	volumeMonts []volumeMountDiff
}

type volumeMountDiff struct {
	diff[k8sv1.VolumeMount]

	data diff[string]
}

func comparePods(pod1, pod2 k8sv1.Pod) podDiff {

}
