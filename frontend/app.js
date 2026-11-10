const lista = document.getElementById("lista-clientes");
const formulario = document.getElementById("form-cliente");
const inputNombre = document.getElementById("nombre");
const mensaje = document.getElementById("mensaje");

// Muestra un error al usuario (y lo deja en consola).
function mostrarError(error) {
  console.error(error);
  mensaje.textContent = error.message;
}

// GET /clientes -> JSON -> elementos del DOM
function cargarClientes() {
  return fetch("/clientes")
    .then((response) => {
      if (!response.ok) {
        throw new Error("No se pudo obtener la lista de clientes");
      }
      return response.json();
    })
    .then((clientes) => {
      lista.innerHTML = "";

      // La API devuelve null (no []) cuando no hay clientes
      if (!clientes || clientes.length === 0) {
        const vacio = document.createElement("li");
        vacio.textContent = "No hay clientes.";
        lista.appendChild(vacio);
        return;
      }

      clientes.forEach((cliente) => {
        const item = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = cliente.nombre; // textContent: no interpreta HTML
        item.appendChild(texto);

        const boton = document.createElement("button");
        boton.textContent = "Eliminar";
        boton.classList.add("eliminar");
        boton.dataset.id = cliente.id_cliente;
        item.appendChild(boton);

        lista.appendChild(item);
      });
    })
    .catch(mostrarError);
}

// POST /clientes al enviar el formulario
formulario.addEventListener("submit", (e) => {
  e.preventDefault();
  mensaje.textContent = "";

  const nombre = inputNombre.value;

  fetch("/clientes", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ nombre }),
  })
    .then((response) => {
      if (!response.ok) {
        // La API responde {"error": "..."} cuando algo falla
        return response.json().then((cuerpo) => {
          throw new Error(cuerpo.error || "No se pudo crear el cliente");
        });
      }
      formulario.reset();
      return cargarClientes();
    })
    .catch(mostrarError);
});

// DELETE /clientes/{id}: delegación de eventos sobre la lista
lista.addEventListener("click", (e) => {
  const boton = e.target.closest("button.eliminar");
  if (!boton) return;

  mensaje.textContent = "";
  const id = boton.dataset.id;

  fetch(`/clientes/${id}`, { method: "DELETE" })
    .then((response) => {
      if (!response.ok) {
        throw new Error("No se pudo eliminar el cliente");
      }
      // 204 No Content: no hay cuerpo que leer
      return cargarClientes();
    })
    .catch(mostrarError);
});

document.addEventListener("DOMContentLoaded", cargarClientes);
