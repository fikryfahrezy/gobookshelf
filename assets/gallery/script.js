const dialogBtn = document.getElementById('dialog-button');
const dialog = document.getElementById('dialog');
const form = document.getElementById('form');
const cancelDialogBtn = document.getElementById('cancel-dialog');

const toogleDialog = function toogleDialog() {
  dialog.classList.toggle('none');
};

if (dialogBtn)
  dialogBtn.addEventListener('click', () => {
    toogleDialog();
  });

cancelDialogBtn.addEventListener('click', () => {
  toogleDialog();
});

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const form = new FormData(e.target);

  fetch('/galleries', {
    method: 'POST',
    body: form,
  })
    .then((res) => {
      if (res.status >= 200 && res.status < 300) {
        window.location.reload();
      }
    })
    .catch((err) => {
      console.log(err);
    });
});
