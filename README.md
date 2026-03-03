Readme realizado con chatgpt

## 1. Validación en servidor (25 puntos)

Se agregaron validaciones directamente en el backend para evitar datos inválidos o inconsistencias.

Las validaciones implementadas incluyen:

* No permitir que el total de episodios sea menor o igual a cero.
* No permitir que los episodios actuales superen el total.
* Validar que los campos requeridos no estén vacíos.
* Evitar valores negativos al disminuir episodios.

Estas validaciones se realizan antes de modificar el slice que almacena las series. Si los datos no cumplen las condiciones, simplemente no se realiza la operación.

La lógica se colocó en los handlers correspondientes (POST, PUT y acciones de actualización).

## 2. Eliminar serie (DELETE) – 20 puntos

Se implementó la funcionalidad para eliminar una serie usando el método DELETE.

Cómo funciona:

* Se agregó un botón "Eliminar" en cada fila de la tabla.
* El botón llama a una función JavaScript que envía una petición DELETE al servidor con el ID de la serie.
* En el backend se busca la serie por ID y se elimina del slice.
* Luego se vuelve a renderizar la tabla actualizada.

El cambio principal fue agregar el botón dentro del loop que genera las filas dinámicamente y crear el handler correspondiente en el servidor.

## 3. Editar serie (PUT) – 25 puntos

Se implementó la posibilidad de editar una serie existente usando el método PUT.

Esto permite modificar información como:

* Nombre
* Total de episodios
* Episodios actuales

Funcionamiento general:

* Se agregó una opción para editar (por ejemplo, mediante un formulario o inputs editables).
* Se envía una petición PUT con los nuevos datos.
* El servidor localiza la serie por ID.
* Se actualizan sus campos después de pasar las validaciones.

La lógica se implementó en un handler específico para PUT, asegurando que no se rompan las reglas de validación (por ejemplo, no permitir que los episodios actuales superen el nuevo total).

## 4. Barra de progreso (15 puntos)

Se agregó una barra visual que representa el porcentaje de episodios vistos.

Cómo se hizo:

* En el backend se calcula el porcentaje usando la relación entre episodios actuales y totales.
* Se genera un bloque HTML con estilos inline que representa una barra.
* Esa barra se inserta como una columna adicional en la tabla.

No se utilizaron librerías externas; la barra está hecha únicamente con HTML y estilos básicos.

## 5. Texto especial si está completa (10 puntos)

Se implementó una condición para mostrar un texto especial cuando una serie está completamente vista.

Cuando:

episodios actuales == total de episodios

Entonces:

* Se muestra un texto como "Completada" o un estado especial junto al nombre o en una columna adicional.

Esta condición se evalúa dentro del loop que genera la tabla. Si se cumple, se concatena el texto especial al momento de construir la fila.

## 6. Botón -1 (10 puntos)

Se agregó un botón "-1" para disminuir el número de episodios vistos.

Funcionamiento:

* El botón llama a una función JavaScript enviando el ID.
* Se hace una petición al backend.
* El servidor reduce el número de episodios actuales en 1.
* Se valida que el valor no sea menor que 0.
* Se actualiza la tabla.

El cambio se hizo agregando el botón dentro del mismo `fmt.Sprintf` que genera cada fila, similar al botón "+1".

# Cambios generales realizados

Los cambios principales se concentraron en:

* Modificar el loop donde se generan las filas dinámicamente.
* Agregar nuevos botones dentro de la tabla.
* Implementar handlers adicionales para DELETE y PUT.
* Añadir validaciones en el backend.
* Agregar lógica condicional para el texto de serie completada.
* Calcular e insertar la barra de progreso.

<img width="1916" height="1079" alt="image" src="https://github.com/user-attachments/assets/f33bcad2-0300-4a6e-9e79-dfd87ee9a4ae" />
