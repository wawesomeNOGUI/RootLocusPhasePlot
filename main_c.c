#ifndef UNICODE
#define UNICODE
#endif 

#include <windows.h>
#include <stdio.h>
#include <stdbool.h>
#include <time.h>

// character(s) to draw on screen
const char charString[] = "GAMER";

// length of charString will be set by strlen(charString) in main
const int CHARS = strlen(charString); 

const COLORREF TRANSPARENT_COLOR = RGB(0, 0, 0);
const COLORREF BACKGROUND_COLOR = RGB(1, 1, 1);

int myWidth, myHeight;
int monitorWidth, monitorHeight;
bool readyToDraw = false;
HWND hwnd;
HDC hDesktopDC, hMyDC, hdcMemDC;
// HRGN hMyRg;
HBITMAP hMyBmp;
RGBQUAD *pPixels;
BITMAPINFO bmi;

LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);
void Wineventproc(
  HWINEVENTHOOK hWinEventHook,
  DWORD event,
  HWND hwnd,
  LONG idObject,
  LONG idChild,
  DWORD idEventThread,
  DWORD dwmsEventTime
);
bool StretchBltToMemDC();
int GetPixelDataFromMemDC();
bool BitBltToWindowDC();
void DrawCharacter(wchar_t *c, int x, int y, int rgb);
void DrawPixel(int x, int y, int rgb);
RGBQUAD *GetColorFromBuffer(RGBQUAD *b, int index);

// call init before any of the other functions in this file
int init()
{
    srand(time(0));
    HINSTANCE hInstance;

    // Register the window class.
    const wchar_t CLASS_NAME[]  = L"Sample Window Class";
    
    WNDCLASS wc = { };

    wc.lpfnWndProc   = WindowProc;
    wc.hInstance     = hInstance;
    wc.lpszClassName = CLASS_NAME;

    RegisterClass(&wc);

    HMONITOR hmon = MonitorFromWindow(GetForegroundWindow(),
                        MONITOR_DEFAULTTONEAREST);

    myWidth = 500;
    myHeight = 500;

    hwnd =  CreateWindowEx(
                0,
                CLASS_NAME,
                L"ASCII Screen",
                WS_SYSMENU | WS_MINIMIZEBOX | WS_VISIBLE,
                0,
                0,
                myWidth,
                myHeight,
                NULL,       // Parent window    
                NULL,       // Menu
                hInstance,  // Instance handle
                NULL        // Additional application data
            );

    // Setup for drawing
    // https://learn.microsoft.com/en-us/windows/win32/gdi/capturing-an-image
	hMyDC = GetDC(hwnd);
    hDesktopDC = GetDC(NULL);

    // Create a compatible DC, which is used in a BitBlt from the window DC.
    hdcMemDC = CreateCompatibleDC(hDesktopDC);
	hMyBmp = CreateCompatibleBitmap(hDesktopDC, myWidth, myHeight);
    SelectObject(hdcMemDC, hMyBmp);

    //set Background Color
    SelectObject(hdcMemDC, CreateSolidBrush(BACKGROUND_COLOR));

    // // This is the best stretch mode.
    // SetStretchBltMode(hMyDC, HALFTONE);
    // SetStretchBltMode(hdcMemDC, HALFTONE);

	// bmi.bmiHeader.biSize = sizeof(bmi.bmiHeader);
	// bmi.bmiHeader.biWidth = myWidth;
	// bmi.bmiHeader.biHeight = -myHeight;  // negative sets origin in top left
	// bmi.bmiHeader.biPlanes = 1;
	// bmi.bmiHeader.biBitCount = 32;
	// bmi.bmiHeader.biCompression = BI_RGB;

	// pPixels = (RGBQUAD*) malloc(myWidth * myHeight * sizeof(RGBQUAD));

	readyToDraw = true;

    // Run the message and update loop.
    // https://learn.microsoft.com/en-us/windows/win32/learnwin32/window-messages
    MSG msg = { };
    while (GetMessage(&msg, NULL, 0, 0) > 0)
    {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    return 0;
}

LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
    switch (uMsg)
    {
    case WM_DESTROY:
        {
            printf("Goodbye!");
            PostQuitMessage(0);
        }
        return 0;
    }
    return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

int GetPixelDataFromMemDC()
{
	// get pixel data into pPixels
    return GetDIBits(
        hdcMemDC,
        hMyBmp,
        0,
        myHeight,
        pPixels,
        &bmi,
        DIB_RGB_COLORS
    );
}

bool BitBltToWindowDC()
{
	// Bit block transfer memdc onto window dc
    return BitBlt(hMyDC,
        0, 0,
        myWidth, myHeight,
        hdcMemDC,
        0, 0,
        SRCCOPY);
}

void DrawCharacter(wchar_t *c, int x, int y, int rgb)
{
	SetTextColor(hdcMemDC, TRANSPARENT_COLOR);
	ExtTextOutW(hdcMemDC, x, y, ETO_OPAQUE, NULL, c, 1, NULL);
}

// rgb integer should have form 0xFFFFFF
void DrawPixel(int x, int y, int rgb)
{
    // draw to buffer DC and when done bit blt to window
    SetPixel(
        hdcMemDC,
        x,
        y,
        RGB((rgb & 0xFF0000) >> 16, (rgb & 0x00FF00) >> 8, rgb & 0x0000FF)
    );
}

RGBQUAD *GetColorFromBuffer(RGBQUAD *b, int index)
{
	return &b[index];
}