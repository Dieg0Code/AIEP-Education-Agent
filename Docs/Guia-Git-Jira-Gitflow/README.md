# 🚀 Guía práctica: Git + Jira + GitFlow

**Profesor:** Diego Obando  
**Para:** Estudiantes del proyecto AIEP

Esta guía establece las convenciones y mejores prácticas para el manejo de código usando Git, integración con Jira y el flujo GitFlow.

## 📋 Resumen ejecutivo

- **Ramas principales:** `main` (producción), `develop` (integración)
- **Ramas de trabajo:** `feature/SCRUM-<n>-descripcion`, `hotfix/SCRUM-<n>-descripcion`, `release/x.y.z`
- **Convención:** Todos los commits y PRs deben incluir la clave de Jira
- **Flujo:** Siempre trabajar desde `develop`, crear PR hacia `develop`

## 🔄 Flujo de trabajo completo

### 1. Preparación

```bash
# Actualizar develop y crear nueva rama
git checkout develop
git pull origin develop
git checkout -b feature/SCRUM-123-descripcion-corta
```

### 2. Desarrollo

```bash
# Commits frecuentes y descriptivos
git add .
git commit -m "SCRUM-123: feat(modulo): implementa funcionalidad X"
git push -u origin feature/SCRUM-123-descripcion-corta
```

### 3. Antes del Pull Request

```bash
# Sincronizar con develop
git fetch origin
git rebase origin/develop  # o git merge origin/develop
# Resolver conflictos si existen
git push --force-with-lease  # solo si hiciste rebase
```

### 4. Pull Request

- **Crear PR** desde tu rama hacia `develop`
- **Título:** `SCRUM-123: Descripción breve de la funcionalidad`
- **Completar plantilla** con enlace a Jira y checklist
- **Solicitar revisión** de 1-2 compañeros

### 5. Merge y cierre

- Esperar aprobaciones y checks verdes
- Mergear el PR
- Actualizar estado en Jira

## 📝 Convenciones obligatorias

### Nombres de ramas

```
feature/SCRUM-123-login-component
hotfix/SCRUM-456-fix-authentication
release/1.2.0
```

### Mensajes de commit

```
SCRUM-123: feat(auth): agrega validación de usuario
SCRUM-123: fix(ui): corrige alineación en formulario
SCRUM-123: docs(readme): actualiza instrucciones de setup
```

**Prefijos recomendados:** `feat`, `fix`, `chore`, `docs`, `test`, `refactor`

### Pull Requests

- **Título:** Debe iniciar con clave Jira
- **Descripción:** Incluir enlace al ticket y checklist completo
- **Target:** Siempre hacia `develop` (excepto hotfixes)

## ✅ Checklist para Pull Requests

Copia este checklist en cada PR:

```markdown
## Checklist

- [ ] Rama creada desde `develop` actualizado
- [ ] Título incluye clave Jira (SCRUM-XXX)
- [ ] Enlace al ticket de Jira en la descripción
- [ ] Código compila sin errores
- [ ] Pruebas locales ejecutadas y exitosas
- [ ] Rebase/merge con `develop` realizado
- [ ] No hay archivos innecesarios (builds, logs, etc.)
- [ ] Revisores asignados (@usuario1, @usuario2)

## Información adicional

- **Ticket Jira:** [SCRUM-XXX](enlace-al-ticket)
- **Tipo:** Feature/Bugfix/Hotfix
- **Screenshots:** (si aplica)
```

## ⚔️ Resolución de conflictos

### Identificar conflictos

```bash
git status  # Ver archivos en conflicto
```

### Resolver manualmente

1. Abrir archivo en VS Code
2. Buscar marcadores: `<<<<<<<`, `=======`, `>>>>>>>`
3. Decidir qué código mantener
4. Eliminar marcadores
5. Guardar archivo

### Completar resolución

```bash
git add archivo-resuelto.js
git rebase --continue  # si estás en rebase
# o
git commit  # si estás en merge
```

### Abortar si es necesario

```bash
git rebase --abort  # volver al estado inicial
git merge --abort   # cancelar merge
```

## 🔧 Comandos de rescate

### Deshacer cambios locales

```bash
git restore archivo.js              # descartar cambios no staged
git restore --staged archivo.js     # quitar archivo del stage
git reset --soft HEAD~1            # deshacer último commit (mantener cambios)
```

### Revertir commits compartidos

```bash
git revert abc123  # crear commit que deshace el commit abc123
```

⚠️ **Evitar:** `git reset --hard` sin supervisión

## 🚨 Hotfixes (urgencias en producción)

```bash
# Crear hotfix desde main
git checkout main
git pull origin main
git checkout -b hotfix/SCRUM-456-fix-critico

# Después del fix
# 1. PR hacia main
# 2. Merge main hacia develop
git checkout develop
git pull origin main
```

## 🔄 Rebase vs Merge

| Aspecto         | Rebase                        | Merge                               |
| --------------- | ----------------------------- | ----------------------------------- |
| **Historial**   | Lineal y limpio               | Conserva contexto original          |
| **Cuándo usar** | Antes de PR, ramas personales | Integraciones, trabajo colaborativo |
| **Comando**     | `git rebase origin/develop`   | `git merge origin/develop`          |

**Recomendación:** Usa rebase para mantener tu rama actualizada antes del PR.

## ❌ Errores comunes y prevención

| Error                       | Prevención                                 |
| --------------------------- | ------------------------------------------ |
| Trabajar en rama incorrecta | Siempre verificar con `git branch`         |
| Push sin actualizar         | `git pull origin develop` antes de empezar |
| Commits gigantes            | Commits pequeños y atómicos                |
| Olvidar clave Jira          | Revisar mensaje antes de commit            |
| Conflictos complejos        | Mantener rama sincronizada frecuentemente  |

## 🎯 Plantillas rápidas

### Crear rama

```bash
git checkout -b feature/SCRUM-[NÚMERO]-[descripcion-corta]
```

### Commit típico

```bash
git commit -m "SCRUM-[NÚMERO]: [tipo]([scope]): [descripción]"
```

### Título de PR

```
SCRUM-[NÚMERO]: [Descripción clara de la funcionalidad]
```

## 📚 Recursos útiles

- **GitFlow oficial:** [Atlassian GitFlow](https://www.atlassian.com/es/git/tutorials/comparing-workflows/gitflow-workflow)
- **Git fundamentals:** [Pro Git Book](https://git-scm.com/book/es/v2)
- **Resolución de conflictos:** [Git Merge Conflicts](https://www.atlassian.com/git/tutorials/using-branches/merge-conflicts)
