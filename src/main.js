import * as tus from 'https://cdn.jsdelivr.net/npm/tus-js-client/+esm';

const input = document.getElementById('fileInput');
const progressBar = document.getElementById('progressBar');
const progressPercent = document.getElementById('progressPercent');
const progressContainer = document.getElementById('progressContainer');

input.addEventListener('change', function (e) {
  progressContainer.style.display = 'block';

  // Get the selected file from the input element
  const file = e.target.files[0]

  // Create a new tus upload
  const upload = new tus.Upload(file, {
    // endpoint URL for the tusd server
    endpoint: 'http://localhost:1337/files/',
    // Retry delays will enable tus-js-client to automatically retry on errors
    retryDelays: [0, 3000, 5000, 10000, 20000],
    // Attach additional meta data about the file for the server
    metadata: {
      filename: file.name,
      filetype: file.type,
    },
    // Callback for errors which cannot be fixed using retries
    onError: function (error) {
      console.log('Failed because: ' + error)
    },
    // Callback for reporting upload progress
    onProgress: function (bytesUploaded, bytesTotal) {
      const percentage = ((bytesUploaded / bytesTotal) * 100).toFixed(2)
      console.log(bytesUploaded, bytesTotal, percentage + '%')

      // progress bar
      progressBar.value  =percentage;
      progressPercent.textContent = percentage + '%';
    },
    // Callback for once the upload is completed
    onSuccess: function () {
      console.log('Download %s from %s', upload.file.name, upload.url)
    },
  })

  // Check if there are any previous uploads to continue.
  upload.findPreviousUploads().then(function (previousUploads) {
    // Found previous uploads so we select the first one.
    if (previousUploads.length) {
      upload.resumeFromPreviousUpload(previousUploads[0])
    }

    // Start the upload
    upload.start()
  })
})
