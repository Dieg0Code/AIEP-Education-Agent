¡Por supuesto, profe Diego! Aquí tienes una guía clara, entretenida y didáctica para tus estudiantes del AIEP sobre **cómo usar Git con Jira y GitFlow**, además de cómo resolver los problemas más comunes. Puedes compartirla como PDF, Notion, Google Doc o imprimirla para el taller.

---

# 🚀 Guía Rápida: Git + Jira + GitFlow

**Metodologías de Desarrollo de Software**  
Profe: Diego Obando

---

## 1. ¿Para qué sirve todo esto? 🤔

- **Git**: Guarda el historial de tu código y permite trabajar en equipo sin pisarse los talones.
- **Jira**: Organiza las tareas del proyecto (como un Trello pro).
- **GitFlow**: Es una forma ordenada de usar ramas en Git para que el caos no reine.

---

## 2. El flujo de trabajo ideal 🏄‍♂️

### 1️⃣ Antes de empezar una tarea

1. **Busca tu tarea en Jira**  
   Ejemplo: `SCRUM-10 Crear pantalla de login`

2. **Crea tu rama feature desde develop**
   ```bash
   git checkout develop
   git pull
   git checkout -b feature/SCRUM-10-login
   ```

---

### 2️⃣ Mientras trabajas

- Haz commits pequeños y frecuentes:
  ```bash
  git add .
  ```

# 🚀 Guía práctica: Git + Jira + GitFlow (mejorada)

Profe: Diego Obando — Versión mejorada: claridad, ejemplos y checklist para PR.

Esta guía está pensada para estudiantes y equipos que participan en el proyecto AIEP. Explica convenciones, comandos útiles, cómo resolver conflictos y buenas prácticas para mantener un repositorio sano.

## Resumen rápido

- Ramas principales: `main` (producción), `develop` (integración)
- Ramas temporales: `feature/*`, `release/*`, `hotfix/*`
- Convención de ramas: `feature/SCRUM-<n>-breve-descripción`
- Convención de commits y PR: incluir la clave de Jira (`SCRUM-10`)

## 1. Flujo recomendado (paso a paso)

1. Antes de empezar

- Revisa la tarea en Jira (ej: `SCRUM-10: Crear pantalla de login`).
- Actualiza `develop` y crea tu rama desde `develop`:

```bash
git checkout develop
git pull origin develop
git checkout -b feature/SCRUM-10-login
```

2. Durante el desarrollo

- Haz commits pequeños y atómicos. Incluye la referencia a Jira en el mensaje.
- Ejemplo de commit:

```bash
git add src/components/Login.tsx
git commit -m "SCRUM-10: feat(login): agrega formulario y validaciones básicas"
git push -u origin feature/SCRUM-10-login
```

3. Antes de crear el Pull Request (PR)

- Asegúrate de traer cambios de `develop` y rebasear o mergearlos localmente:

```bash
git fetch origin
git rebase origin/develop
# o, si prefieres no reescribir historial:
git merge origin/develop
```

- Resuelve conflictos localmente (ver sección de conflictos). Luego push.
- Crea el PR contra `develop`.

4. PR: título y descripción

- Título sugerido: `SCRUM-10: Pantalla de login`
- En la descripción enlaza la historia de Jira y agrega checklist (tests, revisión, screenshots si aplica).

## 2. Convenciones (muy importantes)

- Nombres de ramas:

  - Features: `feature/SCRUM-123-descripcion-corta`
  - Releases: `release/x.y.z`
  - Hotfixes: `hotfix/SCRUM-456-descripcion`

- Mensajes de commit: comienzan con la clave Jira y una etiqueta opcional tipo `feat|fix|chore|docs`.

  - Formato recomendado: `SCRUM-10: feat(login): agrega formulario de login`

- PRs siempre hacia `develop` (salvo hotfixes que vayan a `main` y `develop`).

## 3. Checklist para Pull Requests (pégala en la plantilla de PR)

- [ ] La rama está creada desde `develop`.
- [ ] La descripción incluye la clave Jira y un enlace al ticket.
- [ ] Hice rebase/merge con `develop` y resolví conflictos.
- [ ] Código compilado y pruebas locales pasan.
- [ ] No hay cambios innecesarios en archivos (ej.: package-lock, builds) salvo justificados.
- [ ] Etiqueté a 1-2 revisores.

## 4. Resolución de conflictos (guía práctica)

Cuando haces `rebase` o `merge` y Git te indica conflictos:

1. Identifica archivos en conflicto: Git lista los archivos.
2. Abre en VS Code y busca las marcas `<<<<<<<`, `=======`, `>>>>>>>`.
3. Decide qué quedará (puedes combinar ambos cambios o elegir uno).
4. Marca como resuelto y continúa:

```bash
git add <archivo-resuelto>
git rebase --continue   # si estabas rebaseando
# o
git commit              # si estabas haciendo merge
```

Si necesitas volver al estado previo del rebase:

```bash
git rebase --abort
```

Consejo visual: en VS Code usa la vista de comparación y los botones "Accept Current" / "Accept Incoming".

## 5. Cómo deshacer cambios (operaciones seguras para estudiantes)

- Deshacer cambios en archivos no añadidos:
  - `git restore <archivo>`
- Deshacer cambios staged (antes de commit):
  - `git restore --staged <archivo>`
- Deshacer el último commit pero mantener los cambios en working tree:
  - `git reset --soft HEAD~1`
- Revertir un commit ya compartido (crea un nuevo commit que deshace):
  - `git revert <commit>`

Evita `git reset --hard` sin supervisión: elimina trabajo local.

## 6. Rebase vs Merge (breve)

- Rebase: reescribe historial para aplicar tus commits sobre la punta de `develop`. Mantiene historial lineal.
- Merge: conserva el historial tal cual y crea un commit de merge. No reescribe historial.

Regla sencilla: usa rebase para mantener tu rama actualizada antes del PR (si no has compartido tu rama o el equipo lo acepta). Si no quieres reescribir historial, usa merge.

## 7. Hotfixes y releases (comandos útiles)

- Hotfix (corrección urgente):

```bash
git checkout main
git pull origin main
git checkout -b hotfix/SCRUM-456-fix-login
# arreglas, commiteas, haces PR a main (y luego merge a develop)
```

## 8. Errores comunes y cómo evitarlos

- "Trabajo en la rama equivocada": antes de empezar, comprueba `git status` y la rama con `git branch`.
- "Push sin pull/rebase": siempre trae cambios remotos antes de empezar una tarea.
- "Commits enormes": divide el trabajo en commits lógicos y pequeños.

## 9. Plantillas útiles (copiar/pegar)

- Nombre de rama: `feature/SCRUM-10-login`
- Commit: `SCRUM-10: feat(login): agrega inputs y validación básica`
- PR título: `SCRUM-10: Pantalla de login`
- PR descripción mínima:

```
Relacionado con: SCRUM-10
Descripción: Implementa UI de login, validaciones y rutas.
Checklist:
- [ ] Compila
- [ ] Pruebas unitarias (si aplica)
- [ ] Revisores: @usuario
```

## 10. Recursos

- GitFlow (Atlassian): https://www.atlassian.com/es/git/tutorials/comparing-workflows/gitflow-workflow
- Git basics: https://git-scm.com/book/es/v2
- Resolver conflictos (video): https://www.youtube.com/watch?v=JtIX3HJKwfo

---

Si quieres, puedo:

- Añadir una plantilla de PR en `.github/PULL_REQUEST_TEMPLATE.md`.
- Añadir un archivo `CONTRIBUTING.md` con estas reglas.
- Generar un PDF listo para imprimir.

Dime cuál de estas mejoras quieres que implemente y lo hago.
