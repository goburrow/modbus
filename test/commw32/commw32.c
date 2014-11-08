// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.

// Test serial communication in Win32
// gcc commw32.c
#include <stdlib.h>
#include <stdio.h>
#include <windows.h>

static const char* port = "COM4";

static printLastError() {
	char lpBuffer[256] = "?";
	FormatMessage(FORMAT_MESSAGE_FROM_SYSTEM,
		NULL,
		GetLastError(),
		MAKELANGID(LANG_NEUTRAL,SUBLANG_DEFAULT),
		lpBuffer,
		sizeof(lpBuffer)-1,
		NULL);
	printf("%s\n", lpBuffer);
}

int main(int argc, char** argv) {
	HANDLE handle;
	DCB dcb = {0};
	COMMTIMEOUTS timeouts;
	DWORD n = 0;
	char data[512];
	int i = 0;

	handle = CreateFile(port,
		GENERIC_READ | GENERIC_WRITE,
		0,
		0,
		OPEN_EXISTING,
		0,
		0);
	if (handle == INVALID_HANDLE_VALUE) {
		printf("invalid handle %d\n", GetLastError());
		printLastError();
		return 1;
	}
	printf("handle created %d\n", handle);

	dcb.BaudRate = CBR_9600;
	dcb.ByteSize = 8;
	dcb.StopBits = ONESTOPBIT;
	dcb.Parity = NOPARITY;
	// No software handshaking
	dcb.fTXContinueOnXoff = 1;
	dcb.fOutX = 0;
	dcb.fInX = 0;
	// Binary mode
	dcb.fBinary = 1;
	// No blocking on errors
	dcb.fAbortOnError = 0;

	if (!SetCommState(handle, &dcb)) {
		printf("set comm state error %d\n", GetLastError());
		printLastError();
		CloseHandle(handle);
		return 1;
	}
	printf("set comm state succeed\n");

	// time-out between charactor for receiving (ms)
	timeouts.ReadIntervalTimeout = 1000;
	timeouts.ReadTotalTimeoutMultiplier = 0;
	timeouts.ReadTotalTimeoutConstant = 1000;
	timeouts.WriteTotalTimeoutMultiplier = 0;
	timeouts.WriteTotalTimeoutConstant = 1000;
	if (!SetCommTimeouts(handle, &timeouts)) {
		printf("set comm timeouts error %d\n", GetLastError());
		printLastError();
		CloseHandle(handle);
		return 1;
	}
	printf("set comm timeouts succeed\n");

	if (!WriteFile(handle, "abc", 3, &n, NULL)) {
		printf("write file error %d\n", GetLastError());
		printLastError();
		CloseHandle(handle);
		return 1;
	}
	printf("write file succeed\n");
	printf("Press Enter when ready for reading...");
	getchar();

	if (!ReadFile(handle, data, sizeof(data), &n, NULL)) {
		printf("read file error %d\n", GetLastError());
		printLastError();
		CloseHandle(handle);
		return 1;
	}
	printf("received data %d:\n", n);
	for (i = 0; i < n; ++i) {
		printf("%02x", data[i]);
	}
	printf("\n");

	CloseHandle(handle);
	printf("closed\n");
	return 0;
}
