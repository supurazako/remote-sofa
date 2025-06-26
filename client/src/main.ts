import { Upload, type UploadOptions, type PreviousUpload } from 'tus-js-client';

const input = document.getElementById('fileInput') as HTMLInputElement;
const progressBar = document.getElementById('progressBar') as HTMLProgressElement;
const progressPercent = document.getElementById('progressPercent') as HTMLElement;
const progressContainer = document.getElementById('progressContainer') as HTMLElement;

input.addEventListener('change', (e: Event) => {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) {
    console.error('No file selected');
    return;
  }

  progressContainer.style.display = 'block';

  // Create a new tus upload
  const options: UploadOptions = {
    endpoint: 'http://localhost:1337/files/',
    retryDelays: [0, 3000, 5000, 10000, 20000],
    metadata: {
      filename: file.name,
      filetype: file.type,
    },
    onError: function (error: Error) {
      console.log('Failed because: ' + error)
    },
    onProgress: function (bytesUploaded: number, bytesTotal: number) {
      const percentage = ((bytesUploaded / bytesTotal) * 100).toFixed(2)
      console.log(bytesUploaded, bytesTotal, percentage + '%')

      progressBar.value  = parseFloat(percentage);
      progressPercent.textContent = percentage + '%';
    },
    onSuccess: function () {
      console.log("Success upload");
    },
  };

  // create a new upload instance
  const upload = new Upload(file, options);

  // Check if there are any previous uploads to continue.
  upload.findPreviousUploads().then((previousUploads: PreviousUpload[]) => {
    if (previousUploads.length) {
      upload.resumeFromPreviousUpload(previousUploads[0]);
    }
    upload.start();
  });
});
