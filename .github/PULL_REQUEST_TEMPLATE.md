<!-- Plantilla de Pull Request para el repositorio AIEP-Education-Agent -->
<!-- Instrucciones: al crear un PR, copia/pega y completa los campos. -->

# Título (obligatorio):

# Ejemplo: SCRUM-10: Pantalla de login

## Incluir la clave de Jira (work item key)

Para que Jira muestre este Pull Request en el work item correspondiente debes incluir la clave del trabajo (por ejemplo `JRA-123`, `SCRUM-10`) tanto en el nombre de la rama como en el título del PR.

- Ejemplo de nombre de rama: `git checkout -b SCRUM-10-login`
- Ejemplo de título de PR: `SCRUM-10: Pantalla de login`

GitHub y Jira enlazarán automáticamente el PR al work item cuando detecten la clave en el título (y es buena práctica tenerla también en la rama y en los commits).

## Resumen

Describe brevemente qué hace este PR y por qué. Indica la historia/ticket de Jira relacionada.

- Relacionado con: SCRUM-\_\_\_ (pegue el enlace al ticket de Jira)

- Work item / Jira key: SCRUM-\_\_\_ (ej.: SCRUM-10) - pega el enlace aquí

## Qué se cambió

Lista corta de cambios principales (features, fixes, refactor):

-
-

## Checklist antes de pedir revisión

- [ ] La rama fue creada desde `develop` (o indica si es `main` para hotfix).
- [ ] Hice rebase/merge con `develop` y resolví conflictos localmente.
- [ ] El código compila y la aplicación corre localmente.
- [ ] Las pruebas unitarias/integ. relevantes pasan (o se añadió/actualizó una prueba).
- [ ] No se incluyeron archivos build/lock innecesarios salvo justificación.
- [ ] El mensaje de commit y el título del PR incluyen la clave Jira (ej.: `SCRUM-10`).

- [ ] El nombre de la rama incluye la clave de Jira (ej.: `SCRUM-10-login`).

## Cómo probar los cambios localmente

Incluye pasos mínimos para que un revisor pueda correr tu cambio en su máquina.

```bash
# Ejemplo genérico:
git checkout feature/SCRUM-10-mi-rama
git fetch origin
git rebase origin/develop   # o git merge origin/develop
# build / run commands
```

## Notas sobre seguridad / data / migraciones

Si el PR incluye cambios que afectan a la base de datos, datos sensibles, o configuraciones, detallar aquí.

## Screenshots / evidencia

Adjunta capturas o GIFs si el cambio es visual.

## Revisores sugeridos

Menciona 1-2 personas (ej.: @usuario) o el equipo.

---

Si necesitas ayuda con rebase/conflictos, añade un comentario en el ticket de Jira o etiqueta al profe.
