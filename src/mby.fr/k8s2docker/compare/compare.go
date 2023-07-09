package compare

import (
	"reflect"

	"mby.fr/utils/collections"

	k8sv1 "k8s.io/api/core/v1"
)

type differ interface {
	Path() string
	Left() any
	Right() any
}

type diff[T any] struct {
	path  string
	left  T
	right T
}

func (d diff[T]) Path() string {
	return d.path
}

func (d diff[T]) Left() any {
	return d.left
}

func (d diff[T]) Right() any {
	return d.right
}

type oneSided[T any] struct {
	path string
	item T
	left bool
}

func (d oneSided[T]) Path() string {
	return d.path
}

func (d oneSided[T]) Left() any {
	if d.left {
		return d.item
	}
	return nil
}

func (d oneSided[T]) Right() any {
	if d.left {
		return nil
	}
	return d.item
}

type podDiff struct {
	diffs []differ
}

func (d podDiff) Diffs() []differ {
	return d.diffs
}

func (d podDiff) DiffPathes() []string {
	return collections.Map(d.diffs, func(diff differ) string {
		return diff.Path()
	})
}

func (d podDiff) DoDiffer() bool {
	return len(d.diffs) > 0
}

func (d podDiff) IsUpdatable() bool {
	return len(d.diffs) == len(d.updatableDiffs())
}

func (d podDiff) updatableDiffs() []differ {
	return collections.Filter(d.diffs, func(df differ) bool {
		return df.Path() == "pod.metadata.labels" ||
			df.Path() == "pod.metadata.annotations" ||
			//df.Path() == "pod.spec.hostname" ||
			//df.Path() == "pod.spec.subdomain" ||
			//df.Path() == "pod.spec.hostAliases" ||
			//df.Path() == "pod.spec.dnsConfig" ||
			df.Path() == "pod.spec.restartPolicy" ||
			df.Path() == "pod.spec.volumes" ||
			df.Path() == "pod.spec.containers" ||
			df.Path() == "pod.spec.initContainers" ||
			df.Path() == "pod.spec.containers.livenessProbe" ||
			df.Path() == "pod.spec.containers.readynessProbe" ||
			df.Path() == "pod.spec.containers.startupProbe" ||
			df.Path() == "pod.spec.containers.resources"
	})
}

func ComparePods(left, right k8sv1.Pod) (pd podDiff) {
	d := &pd.diffs
	if anyEquals(left, right) { //  || appendDiff(d, "pod", &left, &right)
		return
	}

	if !anyEquals(left.ObjectMeta, right.ObjectMeta) {
		//if appendDiff(d, "pod.metadata", &left.ObjectMeta, &right.ObjectMeta) {
		appendDiff(d, "pod.metadata.name", &left.ObjectMeta.Name, &right.ObjectMeta.Name)
		appendDiff(d, "pod.metadata.namespace", &left.ObjectMeta.Namespace, &right.ObjectMeta.Namespace)
		appendDiff(d, "pod.metadata.labels", &left.ObjectMeta.Labels, &right.ObjectMeta.Labels)
		appendDiff(d, "pod.metadata.annotations", &left.ObjectMeta.Annotations, &right.ObjectMeta.Annotations)
		//}
	}

	if !anyEquals(left.Spec, right.Spec) {
		volPredicater := func(l, r k8sv1.Volume) bool {
			return l.Name == r.Name
		}
		volDiffAppender := appendDiff[k8sv1.Volume]
		appendArrayDiff(d, "pod.spec.volumes", &left.Spec.Volumes, &right.Spec.Volumes, volPredicater, volDiffAppender)

		ctPredicater := func(l, r k8sv1.Container) bool {
			return l.Name == r.Name
		}
		ctDiffAppender := func(d *[]differ, path string, l, r k8sv1.Container) bool {
			changed := false
			changed = changed || appendDiff(d, path+".image", l.Image, r.Image)
			changed = changed || appendDiff(d, path+".command", l.Command, r.Command)
			changed = changed || appendDiff(d, path+".args", l.Args, r.Args)
			changed = changed || appendDiff(d, path+".workingDir", l.WorkingDir, r.WorkingDir)
			changed = changed || appendDiff(d, path+".ports", l.Ports, r.Ports)
			changed = changed || appendDiff(d, path+".env", l.Env, r.Env)
			changed = changed || appendDiff(d, path+".resources", l.Resources, r.Resources)
			changed = changed || appendDiff(d, path+".volumeMounts", l.VolumeMounts, r.VolumeMounts)
			changed = changed || appendDiff(d, path+".livenessProbe", l.LivenessProbe, r.LivenessProbe)
			changed = changed || appendDiff(d, path+".readinessProbe", l.ReadinessProbe, r.ReadinessProbe)
			changed = changed || appendDiff(d, path+".startupProbe", l.StartupProbe, r.StartupProbe)
			changed = changed || appendDiff(d, path+".imagePullPolicy", l.ImagePullPolicy, r.ImagePullPolicy)
			changed = changed || appendDiff(d, path+".securityContext", l.SecurityContext, r.SecurityContext)
			changed = changed || appendDiff(d, path+".stdin", l.Stdin, r.Stdin)
			changed = changed || appendDiff(d, path+".stdinOnce", l.StdinOnce, r.StdinOnce)
			changed = changed || appendDiff(d, path+".tty", l.TTY, r.TTY)
			return changed
		}
		appendArrayDiff(d, "pod.spec.initContainers", &left.Spec.InitContainers, &right.Spec.InitContainers, ctPredicater, ctDiffAppender)
		appendArrayDiff(d, "pod.spec.containers", &left.Spec.Containers, &right.Spec.Containers, ctPredicater, ctDiffAppender)

		/*
			appendDiff(d, "pod.spec.hostAliases", &left.Spec.HostAliases, &right.Spec.HostAliases)
			appendDiff(d, "pod.spec.restartPolicy", &left.Spec.RestartPolicy, &right.Spec.RestartPolicy)
			appendDiff(d, "pod.spec.securityContext", &left.Spec.SecurityContext, &right.Spec.SecurityContext)
			appendDiff(d, "pod.spec.hostname", &left.Spec.Hostname, &right.Spec.Hostname)
			appendDiff(d, "pod.spec.subdomain", &left.Spec.Subdomain, &right.Spec.Subdomain)
			appendDiff(d, "pod.spec.dnsConfig", &left.Spec.DNSConfig, &right.Spec.DNSConfig)
		*/

		appendDiff(d, "pod.spec.hostAliases", left.Spec.HostAliases, right.Spec.HostAliases)
		appendDiff(d, "pod.spec.restartPolicy", left.Spec.RestartPolicy, right.Spec.RestartPolicy)
		appendDiff(d, "pod.spec.securityContext", left.Spec.SecurityContext, right.Spec.SecurityContext)
		appendDiff(d, "pod.spec.hostname", left.Spec.Hostname, right.Spec.Hostname)
		appendDiff(d, "pod.spec.subdomain", left.Spec.Subdomain, right.Spec.Subdomain)
		appendDiff(d, "pod.spec.dnsConfig", left.Spec.DNSConfig, right.Spec.DNSConfig)
	}

	return pd
}

func appendArrayDiff[T any](d *[]differ, path string, left, right *[]T, predicater func(T, T) bool, diffAppender func(*[]differ, string, T, T) bool) bool {
	ok := false
	if !anyEquals(left, right) {
		var matched []*T
		for _, l := range *left {
			for _, r := range *right {
				if predicater(l, r) {
					ok = ok || diffAppender(d, path, l, r)
					matched = append(matched, &l, &r)
				}
			}
		}

		// Mark removed l
		for _, l := range *left {
			if !collections.ContainsAny[*T](&matched, &l) {
				appendRemovedLeft(d, path, l)
				ok = true
			}
		}

		// Mark added r
		for _, r := range *right {
			if !collections.ContainsAny[*T](&matched, &r) {
				appendNewRight(d, path, r)
				ok = true
			}
		}
	}
	return ok
}

func anyEquals(left, right any) bool {
	return reflect.DeepEqual(left, right)
}

// Add a diff if items are differents
// return true if added a diff
func appendDiff[T any](diffs *[]differ, path string, left, right T) bool {
	//func appendDiff(diffs *[]diff, path string, left, right any) bool {
	/*
		ok := appendNilDiff(diffs, path, left, right)
		if ok {
			return true
		}*/
	/*
		l := reflect.ValueOf(left)
		r := reflect.ValueOf(right)
		if l.Kind() == reflect.Ptr {
			l = reflect.Indirect(l)
			r = reflect.Indirect(r)
			if !reflect.DeepEqual(l.Interface(), r.Interface()) {
				*diffs = append(*diffs, diff[any]{path, l.Interface(), r.Interface()})
				return true
			}
		} else {
			if !reflect.DeepEqual(left, right) {
				*diffs = append(*diffs, diff[T]{path, left, right})
				return true
			}
		}
	*/
	if !reflect.DeepEqual(left, right) {
		*diffs = append(*diffs, diff[T]{path, left, right})
		return true
	}

	return false
}

func appendNewRight[T any](diffs *[]differ, path string, right T) {
	*diffs = append(*diffs, oneSided[T]{path, right, false})
}

func appendRemovedLeft[T any](diffs *[]differ, path string, left T) {
	*diffs = append(*diffs, oneSided[T]{path, left, true})
}
