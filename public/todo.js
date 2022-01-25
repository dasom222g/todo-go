const form = document.querySelector('.form').querySelector('form')
const listArea = document.querySelector('.todo__list')
const header = {
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
}

const getTodos = async () => {
  const response = await fetch('/todos', header)
  const result = await response.json()
  result.forEach((todo) => {
    setTodo(todo)
  })
}

const handleSubmit = async (e) => {
  e.preventDefault()
  const input = form.querySelector('.form__element')
  const title = input.value
  if (/^\s*$/.test(title)) {
    reset(input)
    return
  }
  const data = {
    title,
    is_complete: false
  }
  const response = await fetch('/todos', {
    method: 'POST',
    body: JSON.stringify(data)
  })
  const result = await response.json()
  setTodo(result)
  reset(input)
}

const removeTodo = async (id, element) => {
  const response = await fetch(`todos/${id}`, { method: 'DELETE' })
  const result = await response.json()
  result.success && element.remove()
}

const handleClick = (e) => {
  target = e.target
  switch (target.dataset.eventType) {
    case 'delete':
      element = target.closest('li')
      removeTodo(element.dataset.id, element)
      break
  }
}

// ui
const setTodo = (item) => {
  const li = document.createElement('li')
  li.className = `todo__item`
  li.setAttribute('data-id', item.id)
  const html = `<div class="todo__content ${item.is_complete ? 'complete' : ''}"><div class="todo__item-check"><label><input type="checkbox"><i class="fas fa-square todo__item-check-icon"></i><i class="fas fa-check-square todo__item-check-icon complete"></i><span class="todo__content-text">${item.title}</span></label></div><div class="todo__item-buttonarea"><button type="button" class="todo__item-button"><i class="fas fa-trash-alt" data-event-type="delete"></i></button></div></div>`
  li.innerHTML = html
  item.is_complete && li.querySelector('input').setAttribute('checked', item.is_complete)
  listArea.appendChild(li)
}

const reset = (input) => {
  input.value = ''
  input.focus()
}

getTodos()

listArea.addEventListener('click', handleClick)
form.addEventListener('submit', handleSubmit)