package nsenter

/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

// __attribute__((constructor)) 指的是一旦这个包被引用，该函数就会自动执行
__attribute__((constructor)) void enter_namespace(void) {
	// 从环境变量中获得pid和命令，如果没有指定其中一个就会退出
	char *docker_pid = getenv("docker_pid");
	if (!docker_pid) {
		fprintf(stdout, "missing docker_pid env, skip nsenter");
		return;
	}

	char *docker_cmd = getenv("docker_cmd");
	if (!docker_cmd) {
		fprintf(stdout, "missing docker_cmd env, skip nsenter");
		return;
	}

	int i;
	char nspath[1024];

	// 需要进入的5种namespace
	const unsigned int NAMESPACE_NUM = 5;
	char *namespaces[] = {"ipc", "uts", "net", "pid", "mnt"};
	for (i = 0; i < NAMESPACE_NUM; ++i) {
		// 拼接对应的路径，类似于/proc/pid/ns/ipc
		sprintf(nspath, "/proc/%s/ns/%s", docker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		// 调用setns进入对应的namespace
		if (setns(fd, 0) == -1) {
			fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], stderror(errno));
		} else {
			fprintf(stdout,  "setns on %s namespace succeeded.\n", namespaces[i]);
		}
		close(fd);
	}
	// 在进入的namespace中执行指定的命令
	int res = system(docker_cmd);
	exit(0);
	return;
}
*/
import "C"
