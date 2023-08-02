package driver

var container1_docker_inspect = `
[
	{
		"Id": "83c29efff0c7359784ba8515c20ca77f8f0849bd93ad5da89ca4af5e0e478eb2",
		"Created": "2023-08-01T21:55:46.958074815Z",
		"Path": "ls",
		"Args": [
			"/world"
		],
		"State": {
			"Status": "exited",
			"Running": false,
			"Paused": false,
			"Restarting": false,
			"OOMKilled": false,
			"Dead": false,
			"Pid": 0,
			"ExitCode": 0,
			"Error": "",
			"StartedAt": "2023-08-01T21:55:47.313299229Z",
			"FinishedAt": "2023-08-01T21:55:47.351025854Z"
		},
		"Image": "sha256:a416a98b71e224a31ee99cff8e16063554498227d2b696152a9c3e0aa65e5824",
		"ResolvConfPath": "/var/lib/docker/containers/83c29efff0c7359784ba8515c20ca77f8f0849bd93ad5da89ca4af5e0e478eb2/resolv.conf",
		"HostnamePath": "/var/lib/docker/containers/83c29efff0c7359784ba8515c20ca77f8f0849bd93ad5da89ca4af5e0e478eb2/hostname",
		"HostsPath": "/var/lib/docker/containers/83c29efff0c7359784ba8515c20ca77f8f0849bd93ad5da89ca4af5e0e478eb2/hosts",
		"LogPath": "/var/lib/docker/containers/83c29efff0c7359784ba8515c20ca77f8f0849bd93ad5da89ca4af5e0e478eb2/83c29efff0c7359784ba8515c20ca77f8f0849bd93ad5da89ca4af5e0e478eb2-json.log",
		"Name": "/quizzical_hodgkin",
		"RestartCount": 0,
		"Driver": "overlay2",
		"Platform": "linux",
		"MountLabel": "",
		"ProcessLabel": "",
		"AppArmorProfile": "docker-default",
		"ExecIDs": null,
		"HostConfig": {
			"Binds": [
				"hello:/world"
			],
			"ContainerIDFile": "",
			"LogConfig": {
				"Type": "json-file",
				"Config": {}
			},
			"NetworkMode": "default",
			"PortBindings": {},
			"RestartPolicy": {
				"Name": "no",
				"MaximumRetryCount": 0
			},
			"AutoRemove": false,
			"VolumeDriver": "",
			"VolumesFrom": null,
			"ConsoleSize": [
				39,
				134
			],
			"CapAdd": null,
			"CapDrop": null,
			"CgroupnsMode": "host",
			"Dns": [],
			"DnsOptions": [],
			"DnsSearch": [],
			"ExtraHosts": null,
			"GroupAdd": null,
			"IpcMode": "private",
			"Cgroup": "",
			"Links": null,
			"OomScoreAdj": 0,
			"PidMode": "",
			"Privileged": false,
			"PublishAllPorts": false,
			"ReadonlyRootfs": false,
			"SecurityOpt": null,
			"UTSMode": "",
			"UsernsMode": "",
			"ShmSize": 67108864,
			"Runtime": "runc",
			"Isolation": "",
			"CpuShares": 0,
			"Memory": 0,
			"NanoCpus": 0,
			"CgroupParent": "",
			"BlkioWeight": 0,
			"BlkioWeightDevice": [],
			"BlkioDeviceReadBps": [],
			"BlkioDeviceWriteBps": [],
			"BlkioDeviceReadIOps": [],
			"BlkioDeviceWriteIOps": [],
			"CpuPeriod": 0,
			"CpuQuota": 0,
			"CpuRealtimePeriod": 0,
			"CpuRealtimeRuntime": 0,
			"CpusetCpus": "",
			"CpusetMems": "",
			"Devices": [],
			"DeviceCgroupRules": null,
			"DeviceRequests": null,
			"MemoryReservation": 0,
			"MemorySwap": 0,
			"MemorySwappiness": null,
			"OomKillDisable": false,
			"PidsLimit": null,
			"Ulimits": null,
			"CpuCount": 0,
			"CpuPercent": 0,
			"IOMaximumIOps": 0,
			"IOMaximumBandwidth": 0,
			"MaskedPaths": [
				"/proc/asound",
				"/proc/acpi",
				"/proc/kcore",
				"/proc/keys",
				"/proc/latency_stats",
				"/proc/timer_list",
				"/proc/timer_stats",
				"/proc/sched_debug",
				"/proc/scsi",
				"/sys/firmware"
			],
			"ReadonlyPaths": [
				"/proc/bus",
				"/proc/fs",
				"/proc/irq",
				"/proc/sys",
				"/proc/sysrq-trigger"
			]
		},
		"GraphDriver": {
			"Data": {
				"LowerDir": "/var/lib/docker/overlay2/a62de5840eab18e2662b32a1e7f67e58f4b167fc16cdad4b662eb99b2be9a81f-init/diff:/var/lib/docker/overlay2/dd43539058c9f933a82fd6f3ce3fac76d175f9394fd75e6e3ef910244b454035/diff",
				"MergedDir": "/var/lib/docker/overlay2/a62de5840eab18e2662b32a1e7f67e58f4b167fc16cdad4b662eb99b2be9a81f/merged",
				"UpperDir": "/var/lib/docker/overlay2/a62de5840eab18e2662b32a1e7f67e58f4b167fc16cdad4b662eb99b2be9a81f/diff",
				"WorkDir": "/var/lib/docker/overlay2/a62de5840eab18e2662b32a1e7f67e58f4b167fc16cdad4b662eb99b2be9a81f/work"
			},
			"Name": "overlay2"
		},
		"Mounts": [
			{
				"Type": "volume",
				"Name": "hello",
				"Source": "/var/lib/docker/volumes/hello/_data",
				"Destination": "/world",
				"Driver": "local",
				"Mode": "z",
				"RW": true,
				"Propagation": ""
			}
		],
		"Config": {
			"Hostname": "83c29efff0c7",
			"Domainname": "",
			"User": "",
			"AttachStdin": false,
			"AttachStdout": false,
			"AttachStderr": false,
			"Tty": false,
			"OpenStdin": false,
			"StdinOnce": false,
			"Env": [
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
			],
			"Cmd": [
				"ls",
				"/world"
			],
			"Image": "busybox",
			"Volumes": null,
			"WorkingDir": "",
			"Entrypoint": null,
			"OnBuild": null,
			"Labels": {}
		},
		"NetworkSettings": {
			"Bridge": "",
			"SandboxID": "a843d3dd8dda08a10d353bea8dc0b6f5d3d5fdfd9885f0b59ce0ae9bb88e8d5d",
			"HairpinMode": false,
			"LinkLocalIPv6Address": "",
			"LinkLocalIPv6PrefixLen": 0,
			"Ports": {},
			"SandboxKey": "/var/run/docker/netns/a843d3dd8dda",
			"SecondaryIPAddresses": null,
			"SecondaryIPv6Addresses": null,
			"EndpointID": "",
			"Gateway": "",
			"GlobalIPv6Address": "",
			"GlobalIPv6PrefixLen": 0,
			"IPAddress": "",
			"IPPrefixLen": 0,
			"IPv6Gateway": "",
			"MacAddress": "",
			"Networks": {
				"bridge": {
					"IPAMConfig": null,
					"Links": null,
					"Aliases": null,
					"NetworkID": "ae7de22159ec07d8800941e94caab261546fe6e0d677e70edc2505fc5745e608",
					"EndpointID": "",
					"Gateway": "",
					"IPAddress": "",
					"IPPrefixLen": 0,
					"IPv6Gateway": "",
					"GlobalIPv6Address": "",
					"GlobalIPv6PrefixLen": 0,
					"MacAddress": "",
					"DriverOpts": null
				}
			}
		}
	}
]		
`
