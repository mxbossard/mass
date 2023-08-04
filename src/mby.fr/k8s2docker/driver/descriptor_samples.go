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
			"Labels": {
				"k8s2docker.mby.fr.descriptor.container": "{
					\"securityContext\":{\"runAsNonRoot\": false},\"ports\":[]
				}"
			}
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

// docker run --add-host foo:127.0.0.1 -p 8000:80 -p 8443:443 -v $PWD:/tmp/foo -v bar:/tmp/bar --cpus=0.2 --memory=64m
//
//	-h foohostname --user 10:12 --ip 10.0.0.2 --restart always --workdir /tmp nginx
var container2_docker_inspect = `
[
    {
        "Id": "74b45958736cf0921b317b70d4ac138a3e975fbc1ac32dca2d7fe820f670dc07",
        "Created": "2023-08-03T20:17:33.242120478Z",
        "Path": "/docker-entrypoint.sh",
        "Args": [
            "nginx",
            "-g",
            "daemon off;"
        ],
        "State": {
            "Status": "restarting",
            "Running": true,
            "Paused": false,
            "Restarting": true,
            "OOMKilled": false,
            "Dead": false,
            "Pid": 0,
            "ExitCode": 1,
            "Error": "",
            "StartedAt": "2023-08-03T20:17:51.855063917Z",
            "FinishedAt": "2023-08-03T20:17:52.394234939Z"
        },
        "Image": "sha256:021283c8eb95be02b23db0de7f609d603553c6714785e7a673c6594a624ffbda",
        "ResolvConfPath": "/var/lib/docker/containers/74b45958736cf0921b317b70d4ac138a3e975fbc1ac32dca2d7fe820f670dc07/resolv.conf",
        "HostnamePath": "/var/lib/docker/containers/74b45958736cf0921b317b70d4ac138a3e975fbc1ac32dca2d7fe820f670dc07/hostname",
        "HostsPath": "/var/lib/docker/containers/74b45958736cf0921b317b70d4ac138a3e975fbc1ac32dca2d7fe820f670dc07/hosts",
        "LogPath": "/var/lib/docker/containers/74b45958736cf0921b317b70d4ac138a3e975fbc1ac32dca2d7fe820f670dc07/74b45958736cf0921b317b70d4ac138a3e975fbc1ac32dca2d7fe820f670dc07-json.log",
        "Name": "/wonderful_faraday",
        "RestartCount": 8,
        "Driver": "overlay2",
        "Platform": "linux",
        "MountLabel": "",
        "ProcessLabel": "",
        "AppArmorProfile": "docker-default",
        "ExecIDs": null,
        "HostConfig": {
            "Binds": [
                "/home/maxbundy:/tmp/foo",
                "bar:/tmp/bar"
            ],
            "ContainerIDFile": "",
            "LogConfig": {
                "Type": "json-file",
                "Config": {}
            },
            "NetworkMode": "default",
            "PortBindings": {
                "443/tcp": [
                    {
                        "HostIp": "",
                        "HostPort": "8443"
                    }
                ],
                "80/tcp": [
                    {
                        "HostIp": "",
                        "HostPort": "8000"
                    }
                ]
            },
            "RestartPolicy": {
                "Name": "always",
                "MaximumRetryCount": 0
            },
            "AutoRemove": false,
            "VolumeDriver": "",
            "VolumesFrom": null,
            "ConsoleSize": [
                57,
                118
            ],
            "CapAdd": null,
            "CapDrop": null,
            "CgroupnsMode": "host",
            "Dns": [],
            "DnsOptions": [],
            "DnsSearch": [],
            "ExtraHosts": [
                "foo:127.0.0.1"
            ],
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
            "Memory": 67108864,
            "NanoCpus": 200000000,
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
            "MemorySwap": -1,
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
                "LowerDir": "/var/lib/docker/overlay2/e5f20aeded63443c8a7660deb9e841cad777f9fa9f0d9ebec131dd3109c99935-init/diff:/var/lib/docker/overlay2/a8e2765ca48cdb300551c8ce8a3c948a07d518b4b79ffcf0486bb1b4b1d1906e/diff:/var/lib/docker/overlay2/df1280629735df9a2b69e57a5c63c55b81dfcca12fd1322dd98cc17f730ffe9e/diff:/var/lib/docker/overlay2/3deebc642b103efe2633e5381243d00418b22d92cb72d1ab1bb6b4a71681eb6d/diff:/var/lib/docker/overlay2/00a2aa506675048f47d2774b8a7cca94d8c9fe54d917c365625e471bf71a6b4c/diff:/var/lib/docker/overlay2/132b8ff1e7abd0a3b106863d9c74b0bd14cbd45ccf33e07bd1df411e27f3a148/diff:/var/lib/docker/overlay2/c1afd65e78f122e8f5f34133272a8f5a54ff9440f18d2da9e0ab1c911d7a3175/diff:/var/lib/docker/overlay2/d743b00c70c34d02d78720fb414f851970f942a71bf6b147b9e8c7997b45bdd6/diff",
                "MergedDir": "/var/lib/docker/overlay2/e5f20aeded63443c8a7660deb9e841cad777f9fa9f0d9ebec131dd3109c99935/merged",
                "UpperDir": "/var/lib/docker/overlay2/e5f20aeded63443c8a7660deb9e841cad777f9fa9f0d9ebec131dd3109c99935/diff",
                "WorkDir": "/var/lib/docker/overlay2/e5f20aeded63443c8a7660deb9e841cad777f9fa9f0d9ebec131dd3109c99935/work"
            },
            "Name": "overlay2"
        },
        "Mounts": [
            {
                "Type": "bind",
                "Source": "/home/maxbundy",
                "Destination": "/tmp/foo",
                "Mode": "",
                "RW": true,
                "Propagation": "rprivate"
            },
            {
                "Type": "volume",
                "Name": "bar",
                "Source": "/var/lib/docker/volumes/bar/_data",
                "Destination": "/tmp/bar",
                "Driver": "local",
                "Mode": "z",
                "RW": true,
                "Propagation": ""
            }
        ],
        "Config": {
            "Hostname": "foohostname",
            "Domainname": "",
            "User": "10:12",
            "AttachStdin": false,
            "AttachStdout": true,
            "AttachStderr": true,
            "ExposedPorts": {
                "443/tcp": {},
                "80/tcp": {}
            },
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
                "NGINX_VERSION=1.25.1",
                "NJS_VERSION=0.7.12",
                "PKG_RELEASE=1~bookworm"
            ],
            "Cmd": [
                "nginx",
                "-g",
                "daemon off;"
            ],
            "Image": "nginx",
            "Volumes": null,
            "WorkingDir": "/tmp",
            "Entrypoint": [
                "/docker-entrypoint.sh"
            ],
            "OnBuild": null,
            "Labels": {
				"k8s2docker.mby.fr.descriptor.container": "{
					\"securityContext\":{\"runAsNonRoot\": false},
					\"ports\":[
						{\"name\":\"https\",\"hostPort\":8443,\"containerPort\":443,\"protocol\":\"TCP\"},
						{\"name\":\"http\",\"hostPort\":8000,\"containerPort\":80,\"protocol\":\"TCP\"}
					],
					\"volumeMounts\":[
						{\"name\":\"foo\",\"mountPath\":\"/tmp/foo\"},
						{\"name\":\"bar\",\"mountPath\":\"/tmp/bar\"}
					]
				}",
                "maintainer": "NGINX Docker Maintainers <docker-maint@nginx.com>"
            },
            "StopSignal": "SIGQUIT"
        },
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "e4a06cd923db87018b7aa91d98c9b0a2213573a1f43749d3cc4f4c1545a75280",
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "Ports": {},
            "SandboxKey": "/var/run/docker/netns/e4a06cd923db",
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
                    "NetworkID": "5b18d405f7bee28c53ec5993bce09bcb440119cd4418c5a4e10dfcbba548534f",
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
