package ps

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestScanln(t *testing.T) {
	var in = ` 0.0  0.0 56919 tmux
0.2  0.0 38558 zsh
0.0  0.0 75110 ps
0.0  0.0 89110 tmux
`
	reader := bytes.NewBuffer([]byte(in))
	var cpu, mem float32
	var pid, comm string
	n, err := fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("1: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("1: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("2: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("2: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("3: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("3: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 4 {
		t.Errorf("4: failed to parse 4 value: %d (%s)", n, err)
	} else if err != nil {
		t.Errorf("4: failed to parse: %s", err)
	}
	n, err = fmt.Fscanln(reader, &cpu, &mem, &pid, &comm)
	if n != 0 {
		t.Errorf("5: failed to parse 0 value: %d", n)
	} else if err != io.EOF {
		t.Errorf("5: failed to parse: %s", err)
	}
}

func TestPS(t *testing.T) {
	stats, err := PS()
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) == 0 {
		t.Fatal("no stats")
	}
}

func TestLinux(t *testing.T) {
	reader := bufio.NewReader(bytes.NewBuffer([]byte(linuxPS)))
	reader.ReadLine()
	stats := make([]ProcessStatus, 0)
	for {
		var status ProcessStatus
		n, err := fmt.Fscanln(reader,
			&status.CPU, &status.Mem, &status.PID, &status.Command)
		if n == 4 {
			stats = append(stats, status)
			if err != nil {
				break
			}
		} else if err != nil {
			t.Errorf("reading %d values: %s", n, err)
		}
	}
	for _, stat := range stats {
		if stat.PID == 3781 && stat.CPU != 99.3 {
			t.Error("Failed to parse 3781")
		}
		if stat.PID == 5614 && stat.CPU != 12.3 {
			t.Error("Failed to parse 5614")
		}
	}
}

var linuxPS = `%CPU %MEM     PID COMMAND
 0.0  0.2       1 systemd
 0.0  0.0       2 kthreadd
 0.0  0.0       3 rcu_gp
 0.0  0.0       4 rcu_par_gp
 0.0  0.0       6 kworker/0:0H-kblockd
 0.0  0.0       8 mm_percpu_wq
 0.0  0.0       9 ksoftirqd/0
 0.0  0.0      10 rcuc/0
 0.0  0.0      11 rcu_preempt
 0.0  0.0      12 rcub/0
 0.0  0.0      13 migration/0
 0.0  0.0      14 idle_inject/0
 0.0  0.0      16 cpuhp/0
 0.0  0.0      17 cpuhp/1
 0.0  0.0      18 idle_inject/1
 0.0  0.0      19 migration/1
 0.0  0.0      20 rcuc/1
 0.0  0.0      21 ksoftirqd/1
 0.0  0.0      23 kworker/1:0H-kblockd
 0.0  0.0      24 cpuhp/2
 0.0  0.0      25 idle_inject/2
 0.0  0.0      26 migration/2
 0.0  0.0      27 rcuc/2
 0.0  0.0      28 ksoftirqd/2
 0.0  0.0      30 kworker/2:0H-kblockd
 0.0  0.0      31 cpuhp/3
 0.0  0.0      32 idle_inject/3
 0.0  0.0      33 migration/3
 0.0  0.0      34 rcuc/3
 0.0  0.0      35 ksoftirqd/3
 0.0  0.0      37 kworker/3:0H-kblockd
 0.0  0.0      38 cpuhp/4
 0.0  0.0      39 idle_inject/4
 0.0  0.0      40 migration/4
 0.0  0.0      41 rcuc/4
 0.0  0.0      42 ksoftirqd/4
 0.0  0.0      44 kworker/4:0H-kblockd
 0.0  0.0      45 cpuhp/5
 0.0  0.0      46 idle_inject/5
 0.0  0.0      47 migration/5
 0.0  0.0      48 rcuc/5
 0.0  0.0      49 ksoftirqd/5
 0.0  0.0      50 kworker/5:0-events
 0.0  0.0      51 kworker/5:0H-kblockd
 0.0  0.0      52 cpuhp/6
 0.0  0.0      53 idle_inject/6
 0.0  0.0      54 migration/6
 0.0  0.0      55 rcuc/6
 0.0  0.0      56 ksoftirqd/6
 0.0  0.0      58 kworker/6:0H-kblockd
 0.0  0.0      59 cpuhp/7
 0.0  0.0      60 idle_inject/7
 0.0  0.0      61 migration/7
 0.0  0.0      62 rcuc/7
 0.0  0.0      63 ksoftirqd/7
 0.0  0.0      65 kworker/7:0H-kblockd
 0.0  0.0      66 kdevtmpfs
 0.0  0.0      67 netns
 0.0  0.0      68 rcu_tasks_kthre
 0.0  0.0      69 kauditd
 0.0  0.0      71 kworker/4:1-events
 0.0  0.0      72 khungtaskd
 0.0  0.0      73 oom_reaper
 0.0  0.0      74 writeback
 0.0  0.0      75 kcompactd0
 0.0  0.0      76 ksmd
 0.0  0.0      77 khugepaged
 0.0  0.0     166 kintegrityd
 0.0  0.0     167 kblockd
 0.0  0.0     168 blkcg_punt_bio
 0.0  0.0     169 edac-poller
 0.0  0.0     170 devfreq_wq
 0.0  0.0     171 watchdogd
 0.0  0.0     172 kswapd0
 0.0  0.0     175 kthrotld
 0.0  0.0     176 acpi_thermal_pm
 0.0  0.0     177 nvme-wq
 0.0  0.0     178 nvme-reset-wq
 0.0  0.0     179 nvme-delete-wq
 0.0  0.0     180 ipv6_addrconf
 0.0  0.0     187 kworker/2:1-events
 0.0  0.0     193 kstrp
 0.0  0.0     199 kworker/u17:0
 0.0  0.0     202 kworker/7:1-events
 0.0  0.0     211 kworker/6:1-memcg_kmem_cache
 0.0  0.0     212 charger_manager
 0.0  0.0     245 ata_sff
 0.0  0.0     246 scsi_eh_0
 0.0  0.0     247 scsi_tmf_0
 0.0  0.0     248 scsi_eh_1
 0.0  0.0     249 scsi_tmf_1
 0.0  0.0     250 scsi_eh_2
 0.0  0.0     251 scsi_tmf_2
 0.0  0.0     252 scsi_eh_3
 0.0  0.0     253 scsi_tmf_3
 0.0  0.0     254 scsi_eh_4
 0.0  0.0     255 scsi_tmf_4
 0.0  0.0     256 scsi_eh_5
 0.0  0.0     257 scsi_tmf_5
 0.0  0.0     265 kworker/5:3-cgroup_destroy
 0.0  0.0     266 scsi_eh_6
 0.0  0.0     267 scsi_tmf_6
 0.0  0.0     268 usb-storage
 0.0  0.0     269 kworker/3:2-events
 0.0  0.0     270 kworker/3:1H-kblockd
 0.0  0.0     272 kworker/2:1H-kblockd
 0.0  0.0     273 kworker/5:1H-kblockd
 0.0  0.0     275 kworker/4:1H-kblockd
 0.0  0.0     276 kworker/6:1H-kblockd
 0.0  0.0     277 kworker/1:1H-kblockd
 0.0  0.0     291 kworker/0:1H-kblockd
 0.0  0.0     303 jbd2/sda1-8
 0.0  0.0     304 ext4-rsv-conver
 0.0  0.0     313 kworker/7:1H-kblockd
 0.0  0.2     331 systemd-journal
 0.0  0.0     338 kworker/0:2-events
 0.0  0.0     339 lvmetad
 0.0  0.0     344 iprt-VBoxWQueue
 0.0  0.0     345 iprt-VBoxTscThr
 0.0  0.2     348 systemd-udevd
 0.0  0.0     351 kworker/2:2-events
 0.0  0.2     354 systemd-network
 0.0  0.0     387 kworker/7:2-events
 0.0  0.0     430 nvkm-disp
 0.0  0.0     431 ttm_swap
 0.0  0.4     436 systemd-resolve
 0.0  0.1     437 systemd-timesyn
 0.0  0.0     440 avahi-daemon
 0.0  0.1     441 dbus-daemon
 0.0  0.4     442 NetworkManager
 0.0  0.2     444 upowerd
 0.0  0.1     452 systemd-logind
 0.0  0.0     455 avahi-daemon
 0.0  0.2     461 cupsd
 0.0  0.1     462 sshd
 0.0  0.1     464 gdm
 0.0  0.2     488 accounts-daemon
 0.0  0.5     493 polkitd
 0.0  0.3     495 colord
 0.0  0.3     518 cups-browsed
 0.0  0.0     546 kworker/6:3-events
 0.0  0.0     561 rtkit-daemon
 0.0  0.1     647 wpa_supplicant
 0.0  0.2     784 gdm-session-wor
 0.0  0.2     788 systemd
 0.0  0.0     789 (sd-pam)
 0.0  0.1     801 gdm-x-session
 0.0  1.5     803 Xorg
 0.0  0.1     810 dbus-daemon
 0.0  0.2     812 dwm
 0.0  0.0     821 .xinitrc
 0.0  0.2     829 pulseaudio
 0.0  0.1     831 gsettings-helpe
 0.0  0.0     837 xautolock
 0.0  0.0     838 xcompmgr
 0.0  0.2     839 dunst
 0.0  0.0     846 zsh
 0.0  0.3     853 st
 0.0  0.0     855 tmux: client
 0.0  0.1     857 tmux: server
 0.0  0.1     858 zsh
 0.0  0.0    1219 kworker/1:0-events
 0.0  0.1    1346 gvfsd
 0.0  0.1    1351 gvfsd-fuse
 0.0  0.1    1360 at-spi-bus-laun
 0.0  0.0    1366 dbus-daemon
 0.0  0.1    1393 at-spi2-registr
 0.0  0.1    1445 gvfsd-metadata
 0.0  0.0    2951 kworker/1:1-mm_percpu_wq
 0.0  0.0    2960 kworker/u16:43-events_freezable_power_
 0.0  0.0    2963 kworker/u16:46-events_power_efficient
 0.0  0.0    2968 kworker/u16:51-events_freezable_power_
 0.0  0.0    2975 kworker/u16:58-events_freezable_power_
 0.0  0.0    2983 kworker/u16:66-events_unbound
 0.0  0.0    2986 irq/32-mei_me
 0.0  0.0    2987 kworker/0:0-events
 0.0  0.0    2988 kworker/3:0-mm_percpu_wq
 0.0  0.0    3027 kworker/4:0
 0.0  0.1    3047 abcde
 0.0  0.0    3208 cddb-tool
 0.0  0.1    3209 wget
 0.2  0.3    3210 st
 0.0  0.0    3211 tmux: client
 0.0  0.1    3212 zsh
 0.0  0.2    3217 st
 0.0  0.0    3218 tmux: client
 0.0  0.1    3219 zsh
 0.0  0.2    3772 st
 0.0  0.0    3773 tmux: client
 0.0  0.1    3774 zsh
 0.0  0.1    3780 stress-ng
99.3  0.1    3781 stress-ng-cpu
99.3  0.1    3782 stress-ng-cpu
99.4  0.1    3783 stress-ng-cpu
99.4  0.1    3784 stress-ng-cpu
 0.7  0.5    5606 vim
12.3  8.9    5614 gopls
 0.0  0.0    6622 sleep
 0.0  0.0    6623 ps`
