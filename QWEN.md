# Firefly - Contexto do Projeto

## Visão Geral

**Firefly** é um jogo 2D e starter kit construído com **Ebitengine** (Go). O projeto possui uma arquitetura modular que separa o engine reutilável (`internal/engine`) da implementação específica do jogo (`internal/game`).

### Propósito Atual

O projeto está em evolução de um boilerplate de jogo específico para um **mini-framework pessoal** para criação rápida de jogos 2D com Go/Ebiten. A filosofia é "prático sobre perfeito" - focar em funcionalidades que ajudam a iniciar novos jogos rapidamente.

### Stack Tecnológico

- **Linguagem**: Go 1.25
- **Engine**: Ebitengine v2.8.8
- **UI**: EbitenUI v0.7.2
- **Câmera**: Kamera/v2 v2.97.2
- **Imagens/Fontes**: golang.org/x/image

## Estrutura do Projeto

```
.
├── assets/              # Assets do jogo (imagens, áudio, fonts, particles, tilemap)
├── internal/
│   ├── engine/          # Componentes reutiláveis do engine (framework)
│   │   ├── app/         # Loop principal, contexto, inicialização
│   │   ├── assets/      # Carregamento e gerenciamento de assets
│   │   ├── audio/       # Reprodução de áudio
│   │   ├── contracts/   # Interfaces/interfaces para componentes
│   │   ├── data/config/ # Configuração do engine
│   │   ├── entity/      # Estruturas base (actors, items)
│   │   ├── event/       # Sistema de eventos
│   │   ├── input/       # Input do usuário
│   │   ├── physics/     # Física (body, movement, collision space)
│   │   ├── render/      # Renderização (camera, sprites, tilemap, particles)
│   │   ├── scene/       # Gerenciamento de cenas e transições
│   │   ├── sequences/   # Sequências scriptadas/cutscenes
│   │   ├── ui/          # Componentes de UI (hud, speech)
│   │   └── utils/       # Utilitários (fp16, timing)
│   └── game/            # Implementação específica do jogo Firefly
│       ├── app/         # Setup e inicialização do jogo
│       ├── entity/      # Entidades concretas (actors, items, obstacles)
│       ├── events/      # Eventos específicos do jogo
│       ├── render/      # Lógica de renderização específica
│       ├── scenes/      # Cenas do jogo (intro, menu, phases, story)
│       └── ui/          # UI específica do jogo
├── scripts/
│   └── test_coverage.sh # Script para gerar relatório de coverage
├── coverage/            # Artefatos de teste (coverage.out, coverage.html)
├── main.go              # Entry point da aplicação
├── go.mod               # Definição do módulo Go
└── README.md            # Documentação principal
```

## Build e Execução

### Pré-requisitos

- Go 1.25+
- Git

### Comandos

```bash
# Rodar o jogo
go run main.go

# Rodar todos os testes
go test ./...

# Rodar testes com verbose
go test ./internal/engine/... -v

# Gerar relatório de coverage
./scripts/test_coverage.sh
# ou manualmente:
go test ./... -covermode=atomic -coverprofile=coverage/coverage.out
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Buscar TODOs/FIXMEs
grep -rn "TODO\|FIXME" internal/engine/

# Contar arquivos de teste
find internal/engine -name "*_test.go" | wc -l
```

## Arquitetura e Padrões

### Separação Engine vs Game

- **`internal/engine`**: Componentes genéricos e reutiláveis. Deve ser agnóstico ao jogo.
- **`internal/game`**: Implementação específica do Firefly. Usa contratos do engine.

### Sistema de Contratos (Interfaces)

O engine usa interfaces para desacoplamento. Principais contratos em `internal/engine/contracts/`:

- `body/`: Drawable, MovableCollidableAlive, Shape, Movable, Collidable
- `navigation/`: Scene, Transition, SceneType
- `animation/`: Animation
- `config/`: Configurator
- `context/`: ContextProvider
- `sequences/`: Command, Sequence

### Gerenciamento de Cenas

```go
// Navegação entre cenas
sceneManager.NavigateTo(sceneType, transition, freshInstance)
sceneManager.NavigateBack(transition)
sceneManager.SwitchTo(scene)
```

### Sistema de Física

- **Corpos**: `body.Body` com posição em fixed-point (fp16)
- **Movimento**: Modelos em `physics/movement/` (ex: platformer)
- **Colisão**: Space em `physics/space/` com detecção e resolução
- **Shapes**: Rect, Circle (definidos em `physics/body/`)

### Sistema de Actors

```go
type ActorEntity interface {
    body.Drawable
    Controllable      // OnMoveLeft, OnMoveRight, BlockMovement
    Stateful          // State, SetState, MovementState
    Damageable        // Hurt(damage)
    body.Ownable      // Ownership
    body.MovableCollidableAlive
    
    Update(space body.BodiesSpace) error
    MovementModel() physicsmovement.MovementModel
    GetCharacter() *Character
}
```

### Configuração

O engine suporta configuração via código em `internal/engine/data/config/`:

```go
cfg := config.Get()
cfg.ScreenWidth    // Largura da tela
cfg.ScreenHeight   // Altura da tela
cfg.Physics        // Configurações de física
```

## Convenções de Desenvolvimento

### Estilo de Código

- **Nomenclatura**: PascalCase para exported, camelCase para privado
- **Interfaces**: Nomes descritivos, preferencialmente terminados em `-able` quando aplicável
- **Pacotes**: Nomes curtos e descritivos
- **Erros**: Retornar erros quando possível, `log.Fatal` apenas para erros irrecuperáveis
- **Evitar `_ = variable`**: Não usar `_ = variavel` para silenciar warnings. Usar blank identifier nos parâmetros: `func (t *T) Method(_ Type) {}`

### Práticas de Teste

- Testes unitários para lógica crítica (física, colisões, estados)
- Testes de regressão ao corrigir bugs
- Coverage não é prioridade alta - focar em prevenir regressões
- Usar `go test -covermode=atomic` para testes paralelos seguros

### Debug

- **F1**: Toggle de debug overlay (mostra configs de física)
- Debug drawing disponível em componentes de física
- Logs via `log` package padrão do Go

## Roadmap e Próximos Passos

### Bugs Críticos (TODO.md)

1. **Gap na colisão horizontal ao andar para esquerda** - HIGH
2. **Animação do sprite ao pular** - MEDIUM
3. **Refatorar `shape` para `Shape`** (exported) - LOW

### Melhorias Planejadas

Ver `FRAMEWORK_ROADMAP.md` e `NEXT_STEPS.md` para detalhes completos.

**Fase 1 (Fundação)**:
- [ ] Sistema de configuração JSON
- [ ] Corrigir bugs críticos
- [ ] Melhorar documentação inline
- [ ] Gerador de template de projeto

**Fase 2 (Features)**:
- [ ] Hot-reload de assets (dev)
- [ ] Sistema de eventos simples (se necessário)
- [ ] Suporte a gamepad (se necessário)

**Fase 3+ (Nice-to-have)**:
- [ ] Sistema de save (JSON)
- [ ] Sistema de partículas
- [ ] Melhorias de animação

## Arquivos de Documentação

| Arquivo | Descrição |
|---------|-----------|
| `README.md` | Visão geral da arquitetura e estrutura |
| `FRAMEWORK_ROADMAP.md` | Plano detalhado de evolução do framework |
| `NEXT_STEPS.md` | Próximas ações imediatas e quick wins |
| `TODO.md` | Lista de tarefas e technical debt |
| `internal/engine/README.md` | Documentação do módulo engine |
| `internal/game/README.md` | Documentação do módulo game |

## Dicas para IA/Assistentes

1. **Sempre verificar** `internal/engine/contracts/` antes de implementar novos componentes
2. **Manter separação** engine (reutilável) vs game (específico)
3. **Usar composição** de interfaces em vez de herança
4. **Fixed-point arithmetic** em `internal/engine/utils/fp16/` para posições de física
5. **Scene lifecycle**: `OnStart()` → `Update()` → `Draw()` → `OnFinish()`
6. **Asset loading**: Usar `embed.FS` para assets embutidos
7. **Eventos**: Verificar `internal/engine/event/` para comunicação entre componentes

## Contato e Recursos

- **Repositório**: Local (desenvolvimento pessoal)
- **Ebitengine Docs**: https://ebitengine.org/en/
- **EbitenUI**: https://github.com/ebitenui/ebitenui
