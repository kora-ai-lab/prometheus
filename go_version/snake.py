import pygame
import sys
import time
import random

# Direction constants
UP = 1
RIGHT = 2
DOWN = 3
LEFT = 4

class SnakeGame:
    def __init__(self, width=800, height=600):
        self.width = width
        self.height = height
        self.snake = [(200, 200), (220, 200), (240, 200)]
        self.direction = RIGHT
        self.apple = self.set_new_apple()
        self.score = 0
        self.speed = 10
        self.pause = False
        self.game_over = False
        pygame.init()
        self.display = pygame.display.set_mode((self.width, self.height))
        pygame.display.set_caption('Snake Game')
        self.font = pygame.font.Font(None, 36)
        self.clock = pygame.time.Clock()

    def set_new_apple(self):
        return (random.randint(0, self.width - 20) // 20 * 20,
                random.randint(0, self.height - 20) // 20 * 20)

    def run(self):
        while True:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    pygame.quit()
                    sys.exit()
                elif event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_UP and self.direction != DOWN:
                        self.direction = UP
                    elif event.key == pygame.K_DOWN and self.direction != UP:
                    #    self.direction = DOWN
                    elif event.key == pygame.K_LEFT and self.direction != RIGHT:
                        self.direction = LEFT
                    elif event.key == pygame.K_RIGHT and self.direction != LEFT:
                        self.direction = RIGHT
                    elif event.key == pygame.K_SPACE:
                        self.pause = not self.pause

            if not self.pause and not self.game_over:
                self.move_snake()
                self.check_collision()
                self.update_display()

            self.clock.tick(self.speed)

    def move_snake(self):
        head = self.snake[-1]
        if self.direction == UP:
            new_head = (head[0], head[1] - 20)
        elif self.direction == DOWN:
            new_head = (head[0], head[1] + 20)
        elif self.direction == LEFT:
            new_head = (head[0] - 20, head[1])
        elif self.direction == RIGHT:
            new_head = (head[0] + 20, head[1])

        self.snake.append(new_head)
        if self.apple == new_head:
            self.apple = self.set_new_apple()
            self.score += 1
        else:
            self.snake.pop(0)

    def check_collision(self):
        if (self.snake[-1][0] < 0 or self.snake[-1][0] >= self.width or
            self.snake[-1][1] < 0 or self.snake[-1][1] >= self.height or
            self.snake[-1] in self.snake[:-1]):
            self.game_over = True

    def update_display(self):
        self.display.fill((0, 0, 0))
        for pos in self.snake:
            pygame.draw.rect(self.display, (0, 255, 0), pygame.Rect(pos[0], pos[1], 20, 20))
        pygame.draw.rect(self.display, (255, 0, 0), pygame.Rect(self.apple[0], self.apple[1], 20, 20))
        text = self.font.render(f'Score: {self.score}', True, (255, 255, 255))
        self.display.blit(text, (10, 10))
        pygame.display.flip()

if __name__ == '__main__':
    game = SnakeGame()
    game.run()