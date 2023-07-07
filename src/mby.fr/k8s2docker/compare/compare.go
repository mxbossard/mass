package compare

import (
	"reflect"

	"mby.fr/utils/collections"

	k8sv1 "k8s.io/api/core/v1"
)

type diff struct {
	Path  string
	Left  any
	Right any
}

type podDiff struct {
	diffs []diff
}

func (d podDiff) Diffs() []diff {
	return d.diffs
}

func (d podDiff) DoDiffer() bool {
	return len(d.diffs) > 0
}

func (d podDiff) IsUpdatable() bool {
	return len(d.diffs) == len(d.updatableDiffs())
}

func (d podDiff) updatableDiffs() []diff {
	return collections.Filter(d.diffs, func(df diff) bool {
		return df.Path == "pod.metadata.labels" || df.Path == "pod.metadata.annotations" ||
			df.Path == "pod.spec.hostname" || df.Path == "pod.spec.subdomain" ||
			df.Path == "pod.spec.hostAliases" || df.Path == "pod.spec.dnsConfig"
	})
}

func ComparePods(left, right k8sv1.Pod) (pd podDiff) {
	d := &pd.diffs
	if anyEquals(left, right) || appendNilDiff(d, "pod", &left, &right) {
		return
	}

	if !anyEquals(left.ObjectMeta, right.ObjectMeta) {
		if !appendNilDiff(d, "pod.metadata", &left.ObjectMeta, &right.ObjectMeta) {
			appendDiff(d, "pod.metadata.name", &left.ObjectMeta.Name, &right.ObjectMeta.Name)
			appendDiff(d, "pod.metadata.namespace", &left.ObjectMeta.Namespace, &right.ObjectMeta.Namespace)
			appendDiff(d, "pod.metadata.labels", &left.ObjectMeta.Labels, &right.ObjectMeta.Labels)
			appendDiff(d, "pod.metadata.annotations", &left.ObjectMeta.Annotations, &right.ObjectMeta.Annotations)
		}
	}

	if !anyEquals(left.Spec, right.Spec) {
		if !appendNilDiff(d, "pod.spec", &left.Spec, &right.Spec) {
			//appendDiff(d, "pod.spec.volumes", &left.Spec.Volumes, &right.Spec.Volumes)
			/*
				if !anyEquals(left.Spec.Volumes, right.Spec.Volumes) {
					var matched []k8sv1.Volume
					for _, l := range left.Spec.Volumes {
						for _, r := range right.Spec.Volumes {
							if l.name == r.name {
								appendDiff(d, "pod.spec.volumes", l, r)
								matched = append(matched, l, r)
							}
						}
					}

					// Mark removed l
					for _, l := range left.Spec.Volumes {
						if !collections.Contains(matched, l) {
							appendDiff(d, "pod.spec.volumes", l, nil)
						}
					}

					// Mark added r
					for _, r := range right.Spec.Volumes {
						if !collections.Contains(matched, r) {
							appendDiff(d, "pod.spec.volumes", nil, r)
						}
					}
				}
			*/

			volPredicater := func(l, r k8sv1.Volume) bool {
				return l.Name == r.Name
			}
			volDiffAppender := appendDiff[k8sv1.Volume]
			appendArrayDiff(d, "pod.spec.volumes", &left.Spec.Volumes, &right.Spec.Volumes, volPredicater, volDiffAppender)

			appendDiff(d, "pod.spec.initContainers", &left.Spec.InitContainers, &right.Spec.InitContainers)
			appendDiff(d, "pod.spec.containers", &left.Spec.Containers, &right.Spec.Containers)
			appendDiff(d, "pod.spec.restartPolicy", &left.Spec.RestartPolicy, &right.Spec.RestartPolicy)
			appendDiff(d, "pod.spec.securityContext", &left.Spec.SecurityContext, &right.Spec.SecurityContext)
			appendDiff(d, "pod.spec.hostname", &left.Spec.Hostname, &right.Spec.Hostname)
			appendDiff(d, "pod.spec.subdomain", &left.Spec.Subdomain, &right.Spec.Subdomain)
			appendDiff(d, "pod.spec.hostAliases", &left.Spec.HostAliases, &right.Spec.HostAliases)
			appendDiff(d, "pod.spec.dnsConfig", &left.Spec.DNSConfig, &right.Spec.DNSConfig)
		}
	}

	return pd
}

func appendArrayDiff[T any](d *[]diff, path string, left, right *[]T, predicater func(T, T) bool, diffAppender func(*[]diff, string, *T, *T) bool) bool {
	ok := false
	if !anyEquals(left, right) {
		var matched []*T
		for _, l := range *left {
			for _, r := range *right {
				if predicater(l, r) {
					ok = ok || diffAppender(d, path, &l, &r)
					matched = append(matched, &l, &r)
				}
			}
		}

		// Mark removed l
		for _, l := range *left {
			if !collections.Contains(matched, &l) {
				ok = ok || diffAppender(d, path, &l, nil)
			}
		}

		// Mark added r
		for _, r := range *right {
			if !collections.Contains(matched, &r) {
				ok = ok || diffAppender(d, path, nil, &r)
			}
		}
	}
	return ok
}

func anyEquals(left, right any) bool {
	return reflect.DeepEqual(left, right)
}

// Add a diff if one of items is nil
// return true if added a diff
func appendNilDiff[T any](diffs *[]diff, path string, left, right *T) bool {
	if (left == nil || right == nil) && left != right {
		*diffs = append(*diffs, diff{path, left, right})
		return true
	}
	return false
}

// Add a diff if items are differents
// return true if added a diff
func appendDiff[T any](diffs *[]diff, path string, left, right *T) bool {
	ok := appendNilDiff(diffs, path, left, right)
	if ok {
		return true
	}
	if !reflect.DeepEqual(left, right) {
		*diffs = append(*diffs, diff{path, left, right})
		return true
	}
	return false
}
