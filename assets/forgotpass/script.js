const form = document.getElementById('form');
const message = document.getElementById('message');

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const fData = new FormData(e.target);
  const data = Object.fromEntries(fData.entries());

  fetch('/forgotpassword', {
    headers: {
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(data),
  })
    .then((res) => {
      if (!(res.status >= 200 && res.status < 300)) {
        return res.json();
      }
    })
    .then((res) => {
      if (res) throw new Error(res.message);

      message.textContent = 'Success, check your email!';
    })
    .catch((err) => {
      message.textContent = err.message;
    });
});
