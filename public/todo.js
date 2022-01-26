const form = document.querySelector('.form').querySelector('form')
const listArea = document.querySelector('.todo__list')
// const completeCheckbox = listArea.querySelector(`input[type='checkbox']`)
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

const changeCompleteTodo = async (id, checked, element) => {
  // 'todos/1?complete=true'
  const content = element.querySelector('.todo__content')
  const response = await fetch(`todos/${id}?complete=${checked}`)
  const result = await response.json()
  content.classList.remove('complete') // 초기화
  console.log('result', result)
  checked && content.classList.add('complete')
}

const handleClick = (e) => {
  target = e.target
  const element = target.closest('li')
  switch (target.dataset.eventType) {
    case 'delete':
      removeTodo(element.dataset.id, element)
      break
    case 'changeComplete':
      changeCompleteTodo(element.dataset.id, target.checked, element)
      break
  }
}

// ui
const setTodo = (item) => {
  const li = document.createElement('li')
  li.className = `todo__item`
  li.setAttribute('data-id', item.id)
  const html = `<div class="todo__content ${item.is_complete ? 'complete' : ''}"><div class="todo__item-check"><label><input type="checkbox" data-event-type="changeComplete"><i class="fas fa-square todo__item-check-icon"></i><i class="fas fa-check-square todo__item-check-icon complete"></i><span class="todo__content-text">${item.title}</span></label></div><div class="todo__item-buttonarea"><button type="button" class="todo__item-button"><i class="fas fa-trash-alt" data-event-type="delete"></i></button></div></div>`
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