const form = document.getElementById('form');

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const form = new FormData(e.target);
  const data = Object.fromEntries(form.entries());

  fetch('/loginacc', {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(data),
  })
    .then((res) => {
      if (res.status >= 200 && res.status < 300) {
        window.location = '/';
      }
    })
    .catch((err) => {
      console.log(err);
    });
});
