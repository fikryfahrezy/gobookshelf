const form = document.getElementById('form');

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const url = new URL(window.location);
  const sParam = new URLSearchParams(url.search);
  const fData = new FormData(e.target);
  const data = Object.fromEntries(fData.entries());

  data.code = sParam.get('code');

  fetch('/updatepassword', {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'PATCH',
    body: JSON.stringify(data),
  })
    .then((res) => {
      if (res.status >= 200 && res.status < 300) {
        window.location = '/login';
      }
    })
    .catch((err) => {
      console.log(err);
    });
});
