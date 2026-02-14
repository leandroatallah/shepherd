# Firefly

A 2D game built with Ebitengine, featuring a modular architecture.

## Architecture Overview

The project is structured into two main packages: `engine` and `game`.

- **`internal/engine`**: The core game engine, providing reusable components for scenes, physics, actors, and other systems.
- **`internal/game`**: The specific implementation of the game, including scenes, characters, and items.

This separation allows the engine to be developed independently from the game's content.

## Folder Structure

```
.
├── assets/              # Game assets (images, sounds, etc.)
│   ├── audio/           # Audio files
│   ├── fonts/           # Font files
│   ├── images/          # Image files
│   ├── particles/       # Particle effect configurations
│   └── tilemap/         # Tilemap related assets
├── main.go            # Application entry point
├── internal/
│   ├── engine/          # Core game engine components
│   │   ├── app/         # Main engine loop, context, and initialization
│   │   ├── assets/      # Asset loading and management (images, fonts)
│   │   ├── audio/       # Audio playback functionality
│   │   ├── contracts/   # Interfaces for engine components (animation, body, config, context, navigation, sequences, tilemaplayer)
│   │   ├── data/        # Data loading, management, and configuration schemas
│   │   │   └── config/  # Engine configuration
│   │   ├── entity/      # Foundational structures for in-game objects (actors, items)
│   │   │   ├── actors/  # Actor management and movement (e.g., characters, enemies)
│   │   │   └── items/   # Item management
│   │   ├── event/       # Event handling system
│   │   ├── input/       # User input handling
│   │   ├── physics/     # Physics simulation (body, movement, skill, space)
│   │   ├── render/      # Rendering tasks (camera, particles, screenutil, sprites, tilemap)
│   │   │   └── camera/  # Camera control and rendering
│   │   ├── scene/       # Game scene management and transitions
│   │   ├── sequences/   # Game sequences and command processing
│   │   ├── ui/          # Building blocks for user interface elements (hud, speech)
│   │   └── utils/       # Utility functions
│   └── game/            # Game-specific implementation
│       ├── app/         # Game-specific setup and initialization (config, phases list, setup)
│       ├── entity/      # Concrete game entities (actors, items, obstacles, types)
│       │   ├── actors/  # Game-specific characters and enemies
│       │   ├── items/   # Game-specific items
│       │   └── obstacles/ # Game-specific obstacles
│       ├── events/      # Game-specific events
│       ├── render/      # Game-specific rendering logic
│       │   └── camera/  # Game-specific camera settings
│       ├── scenes/      # Game scenes and phases (intro, menu, phases, story, summary, types)
│       ├── ui/          # Game's specific user interface (hud, speech)
│       └── README.md
├── go.mod               # Go module definition
└── README.md
```

## Dependencies

- **Ebitengine**: A dead simple 2D game engine for Go.
- **EbitenUI**: A UI library for Ebitengine.
- **Kamera/v2**: A camera library for Ebitengine.
- **Go**: The programming language.
