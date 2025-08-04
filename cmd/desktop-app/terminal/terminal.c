#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <termios.h>
#include <sys/ioctl.h>
#include <signal.h>

// C structs matching the Go bridge
typedef struct {
    int x, y;
    int glyph;
    int foregroundR, foregroundG, foregroundB;
    int backgroundR, backgroundG, backgroundB;
    int style;
    float alpha;
} CGlyph;

typedef struct {
    CGlyph* glyphs;
    int width, height;
    int count;
} CGlyphMatrix;

// External functions from Go
extern void initGame(int width, int height);
extern void updateGame(float dt, int thrust, int brake, int left, int right);
extern CGlyphMatrix getGlyphMatrix();

// Terminal handling
struct termios original_termios;
int term_width = 80;
int term_height = 24;

void cleanup() {
    tcsetattr(STDIN_FILENO, TCSAFLUSH, &original_termios);
    printf("\033[?1049l"); // Exit alt screen
    printf("\033[?25h");   // Show cursor
}

void setup_terminal() {
    // Get terminal size
    struct winsize w;
    if (ioctl(STDOUT_FILENO, TIOCGWINSZ, &w) == 0) {
        term_width = w.ws_col;
        term_height = w.ws_row;
    }
    
    // Setup raw mode
    tcgetattr(STDIN_FILENO, &original_termios);
    atexit(cleanup);
    
    struct termios raw = original_termios;
    raw.c_lflag &= ~(ECHO | ICANON | ISIG | IEXTEN);
    raw.c_iflag &= ~(IXON | ICRNL | BRKINT | INPCK | ISTRIP);
    raw.c_cflag |= CS8;
    raw.c_oflag &= ~(OPOST);
    raw.c_cc[VMIN] = 0;
    raw.c_cc[VTIME] = 1;
    
    tcsetattr(STDIN_FILENO, TCSAFLUSH, &raw);
    
    // Alt screen
    printf("\033[?1049h"); // Enter alt screen
    printf("\033[?25l");   // Hide cursor
    printf("\033[2J");     // Clear screen
}

void render_frame(CGlyphMatrix matrix) {
    printf("\033[H"); // Move cursor to home
    
    if (matrix.glyphs == NULL || matrix.count == 0) {
        // Render empty screen
        for (int y = 0; y < term_height; y++) {
            for (int x = 0; x < term_width; x++) {
                printf(".");
            }
            printf("\n");
        }
        return;
    }
    
    // Create screen buffer
    char** screen = malloc(matrix.height * sizeof(char*));
    for (int y = 0; y < matrix.height; y++) {
        screen[y] = malloc((matrix.width + 1) * sizeof(char));
        for (int x = 0; x < matrix.width; x++) {
            screen[y][x] = '.'; // Default background
        }
        screen[y][matrix.width] = '\0';
    }
    
    // Fill with glyphs
    for (int i = 0; i < matrix.count; i++) {
        CGlyph glyph = matrix.glyphs[i];
        if (glyph.y >= 0 && glyph.y < matrix.height && 
            glyph.x >= 0 && glyph.x < matrix.width) {
            screen[glyph.y][glyph.x] = (char)glyph.glyph;
        }
    }
    
    // Render screen with colors
    for (int y = 0; y < matrix.height && y < term_height; y++) {
        for (int x = 0; x < matrix.width && x < term_width; x++) {
            // Find glyph for color info
            char ch = screen[y][x];
            int r = 255, g = 255, b = 255; // Default white
            
            for (int i = 0; i < matrix.count; i++) {
                CGlyph glyph = matrix.glyphs[i];
                if (glyph.x == x && glyph.y == y) {
                    r = glyph.foregroundR;
                    g = glyph.foregroundG;
                    b = glyph.foregroundB;
                    break;
                }
            }
            
            // Output with color
            printf("\033[38;2;%d;%d;%dm%c", r, g, b, ch);
        }
        printf("\033[0m\n"); // Reset color and newline
    }
    
    // Cleanup screen buffer
    for (int y = 0; y < matrix.height; y++) {
        free(screen[y]);
    }
    free(screen);
    
    fflush(stdout);
}

int main() {
    printf("Harvester Desktop Renderer (Terminal Mode)\n");
    printf("W=thrust, S=brake, A/D=turn, Q=quit\n");
    printf("Press any key to start...\n");
    getchar();
    
    setup_terminal();
    
    // Initialize game
    initGame(term_width, term_height);
    
    int thrust = 0, brake = 0, left = 0, right = 0;
    
    printf("Game initialized. Controls: WASD + Q to quit\n");
    sleep(1);
    
    while (1) {
        // Check for input (non-blocking)
        char c;
        if (read(STDIN_FILENO, &c, 1) == 1) {
            switch (c) {
                case 'q': 
                case 'Q': 
                case 27: // ESC
                    goto cleanup_exit;
                case 'w':
                case 'W':
                    thrust = 1;
                    break;
                case 's':
                case 'S':
                    brake = 1;
                    break;
                case 'a':
                case 'A':
                    left = 1;
                    break;
                case 'd':
                case 'D':
                    right = 1;
                    break;
            }
        } else {
            // Reset inputs when no key pressed
            thrust = brake = left = right = 0;
        }
        
        // Update game
        updateGame(0.016, thrust, brake, left, right);
        
        // Get and render frame
        CGlyphMatrix matrix = getGlyphMatrix();
        render_frame(matrix);
        
        // ~60 FPS
        usleep(16667);
    }
    
cleanup_exit:
    cleanup();
    printf("Thanks for playing Harvester!\n");
    return 0;
}