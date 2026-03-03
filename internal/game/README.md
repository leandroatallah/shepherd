# Game Module

This module contains the specific implementation and logic for the _Shepherd_ game. It is built upon the reusable components and contracts provided by the `internal/engine` module.

## Game-Specific Logic

- `app/`: Contains the game-specific setup and initialization code, configuring the engine to run _Shepherd_.
  - `config.go`: Game configuration settings.
  - `phases_list.go`: Defines the list and order of game phases.
  - `setup.go`: Handles game initialization and setup routines.
- `entity/`: Defines the concrete game entities.
  - `actors/`: Implements the `Player` and specific `Enemies` for the game.
  - `items/`: Implements collectible items like `Coin`.
  - `obstacles/`: Defines game-specific obstacles.
  - `types/`: Custom types related to game entities.
- `scenes/`: Implements the actual game scenes, such as the `IntroScene`, `MenuScene`, and gameplay levels. It orchestrates the actors, items, and UI for each part of the game.
  - `init_scenes.go`: Initializes all game scenes.
  - `scene_intro.go`: The introductory scene.
  - `scene_menu.go`: The main menu scene.
  - `scene_phase_title.go`: Scene for displaying the phase title.
  - `scene_phase_reboot.go`: Scene for rebooting a phase.
  - `scene_story.go`: Scenes dedicated to story progression.
  - `scene_summary.go`: Scene for displaying game summary or scores.
  - `sounds.go`: Management of game-specific audio triggers within scenes.
  - `phases/`: Specific implementations for game phases within scenes.
  - `types/`: Custom types for game scenes.

## Customization and Implementation

- `render/`: Contains game-specific rendering logic.
  - `vfx/`: Game-specific visual effects, such as overhead and screen text.
- `ui/`: Implements the game's specific user interface.
  - `hud/`: Game's main Heads-Up Display elements.
  - `speech/`: Game-specific speech bubbles and dialogue styles.
