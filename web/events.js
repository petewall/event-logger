const labelColors = {
  customer: 'blue',
  people: 'orange'
}

function addLabel() {
  const key = $('.label-key').val()
  const value = $('.label-value').val()
 
  const label = $('<a>', {
    class: `ui ${labelColors[key]} label`
  }).append(
    $('<span>', {text: value}),
    $('<div>', {class: 'detail', text: key})
  )
  label.click(() => {
    label.remove()
  })
  $('#label-list').append(label)
}

function showNewModal() {
  $('#label-list').empty()
  $('.ui.new.modal').modal({
    onApprove: () => {
      const time = $('.ui.new.modal input.time').val()
      const message = $('.ui.new.modal input.message').val()
      const labels = {}
      for (let label of $('#label-list').children()) {
        const key = $(label).children('div.detail').text()
        const value = $(label).children('span').text()
        if (labels[key]) {
          labels[key] += `|${value}`
        } else {
          labels[key] = value
        }
      }
      const newEvent = { time, message, labels }
      $.ajax({
        url: '/api/event',
        type: 'PUT',
        data: JSON.stringify(newEvent),
        contentType: 'application/json',
        success: function() {
          location.reload()
        }
      })
    }
  })
  $('.ui.new.modal').modal('show')

}

function setEvents(events) {
  $('#events-list').empty()
  for (const event of events) {
    const labels = $('<td>')

    for (let key in event.labels) {
      labels.append(
        $('<a>', {
          class: `ui ${labelColors[key]} label`
        }).append(
          event.labels[key],
          $('<div>', {class: 'detail', text: key})
        )
      )
    }

    $('#events-list').append(
      $('<tr>').append(
        $('<td>', {text: event.time}),
        $('<td>', {text: event.message}),
        labels
      ),
    )
  }
}

$(document).ready(() => {
  $.get('http://localhost:8081/api/events', (data) => {
    const events = JSON.parse(data)
    setEvents(events)
  })
  $('.new.button').click(showNewModal)
  $('.new.modal #add-label-button').click(addLabel)
})
