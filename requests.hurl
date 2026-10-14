# requests.hurl
#
# Prueba de punta a punta de la API REST del Trabajo de Cursada.
# Requiere el servidor corriendo (ver README.md).
#
# Ejecutar con:
#   hurl --test requests.hurl


# ============================================================
# CLIENTES
# ============================================================

# 1. Crear un cliente
POST http://localhost:8090/clientes
Content-Type: application/json
{
  "nombre": "Acme S.A."
}
HTTP 201
[Asserts]
jsonpath "$.id_cliente" isInteger
jsonpath "$.nombre" == "Acme S.A."
[Captures]
cliente_id: jsonpath "$.id_cliente"

# 2. Validación: rechazar nombre vacío
POST http://localhost:8090/clientes
Content-Type: application/json
{
  "nombre": ""
}
HTTP 400

# 3. Listar clientes
GET http://localhost:8090/clientes
HTTP 200
[Asserts]
jsonpath "$" isCollection

# 4. Obtener el cliente por ID
GET http://localhost:8090/clientes/{{cliente_id}}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Acme S.A."

# 5. Obtener un cliente inexistente -> 404
GET http://localhost:8090/clientes/999999
HTTP 404

# 6. Actualizar el cliente
PUT http://localhost:8090/clientes/{{cliente_id}}
Content-Type: application/json
{
  "nombre": "Acme Corp."
}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Acme Corp."

# 7. Eliminar el cliente
DELETE http://localhost:8090/clientes/{{cliente_id}}
HTTP 204

# 8. Verificar que ya no existe
GET http://localhost:8090/clientes/{{cliente_id}}
HTTP 404


# ============================================================
# DESARROLLADORES
# ============================================================

# 1. Crear un desarrollador
POST http://localhost:8090/desarrolladores
Content-Type: application/json
{
  "nombre": "Ada",
  "apellido": "Lovelace"
}
HTTP 201
[Asserts]
jsonpath "$.id_desarrollador" isInteger
[Captures]
dev_id: jsonpath "$.id_desarrollador"

# 2. Validación: rechazar apellido vacío
POST http://localhost:8090/desarrolladores
Content-Type: application/json
{
  "nombre": "Ada",
  "apellido": ""
}
HTTP 400

# 3. Listar desarrolladores
GET http://localhost:8090/desarrolladores
HTTP 200
[Asserts]
jsonpath "$" isCollection

# 4. Obtener el desarrollador por ID
GET http://localhost:8090/desarrolladores/{{dev_id}}
HTTP 200
[Asserts]
jsonpath "$.apellido" == "Lovelace"

# 5. Actualizar el desarrollador
PUT http://localhost:8090/desarrolladores/{{dev_id}}
Content-Type: application/json
{
  "nombre": "Ada",
  "apellido": "King"
}
HTTP 200
[Asserts]
jsonpath "$.apellido" == "King"

# (No lo eliminamos todavía: se usa en la sección de la relación N:M)


# ============================================================
# PROYECTOS
# ============================================================

# 1. Crear un cliente para asociar al proyecto
POST http://localhost:8090/clientes
Content-Type: application/json
{
  "nombre": "Globex Inc."
}
HTTP 201
[Captures]
proyecto_cliente_id: jsonpath "$.id_cliente"

# 2. Crear un proyecto
POST http://localhost:8090/proyectos
Content-Type: application/json
{
  "nombre": "Sistema de Facturación",
  "descripcion": "Migración del sistema legado",
  "estado": "en_progreso",
  "gastos": 1500.50,
  "presupuesto": 20000.00,
  "tiempo_estimado": "2026-12-31T00:00:00Z",
  "id_cliente": {{proyecto_cliente_id}}
}
HTTP 201
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.nombre" == "Sistema de Facturación"
[Captures]
proyecto_id: jsonpath "$.id"

# 3. Validación: rechazar presupuesto negativo
POST http://localhost:8090/proyectos
Content-Type: application/json
{
  "nombre": "Proyecto inválido",
  "descripcion": "-",
  "estado": "nuevo",
  "gastos": 0,
  "presupuesto": -100,
  "tiempo_estimado": "2026-12-31T00:00:00Z",
  "id_cliente": {{proyecto_cliente_id}}
}
HTTP 400

# 4. Listar proyectos
GET http://localhost:8090/proyectos
HTTP 200
[Asserts]
jsonpath "$" isCollection

# 5. Obtener el proyecto por ID
GET http://localhost:8090/proyectos/{{proyecto_id}}
HTTP 200
[Asserts]
jsonpath "$.estado" == "en_progreso"

# 6. Actualizar el proyecto
PUT http://localhost:8090/proyectos/{{proyecto_id}}
Content-Type: application/json
{
  "nombre": "Sistema de Facturación",
  "descripcion": "Migración del sistema legado",
  "estado": "completado",
  "gastos": 19800.00,
  "presupuesto": 20000.00,
  "tiempo_estimado": "2026-12-31T00:00:00Z",
  "id_cliente": {{proyecto_cliente_id}}
}
HTTP 200
[Asserts]
jsonpath "$.estado" == "completado"


# ============================================================
# RELACIÓN N:M: DESARROLLADORES EN PROYECTOS
# ============================================================

# 1. Asignar el desarrollador creado antes a este proyecto
POST http://localhost:8090/proyectos/{{proyecto_id}}/desarrolladores
Content-Type: application/json
{
  "id_desarrollador": {{dev_id}}
}
HTTP 201
[Asserts]
jsonpath "$.id_desarrollador" == {{dev_id}}
jsonpath "$.id_proyecto" == {{proyecto_id}}

# 2. Listar los desarrolladores del proyecto
GET http://localhost:8090/proyectos/{{proyecto_id}}/desarrolladores
HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[0].id_desarrollador" == {{dev_id}}

# 3. Quitar al desarrollador del proyecto
DELETE http://localhost:8090/proyectos/{{proyecto_id}}/desarrolladores/{{dev_id}}
HTTP 204

# 4. Volver a quitarlo -> ya no existe la asignación -> 404
DELETE http://localhost:8090/proyectos/{{proyecto_id}}/desarrolladores/{{dev_id}}
HTTP 404


# ============================================================
# LIMPIEZA
# ============================================================

DELETE http://localhost:8090/proyectos/{{proyecto_id}}
HTTP 204

DELETE http://localhost:8090/desarrolladores/{{dev_id}}
HTTP 204

DELETE http://localhost:8090/clientes/{{proyecto_cliente_id}}
HTTP 204
