# Form signer

Allow patients to fill in and sign an html form, submit the results and get served a pdf.
The pdf content will contain the html form content, the current date and their
signature.

The golang code is designed to be able to parse different html templates.

## Use

Currently hosted on google app engine at http://dinodent-forms.appspot.com/new-patient.html.

## Development

### Checkout

Checkout using git clone. Remember to first set your GOPATH variable to the project root.

### Google App engine usage reminder

- Download and install [gcloud sdk](https://cloud.google.com/sdk/docs/)
- Authenticate my account (may happen as part of installation):

      gcloud init

- As per the [quick start](https://cloud.google.com/appengine/docs/standard/go/quickstart), after installing app-engine-go you can run a local development server:

      gcloud components install app-engine-go
      dev_appserver.py app.yaml

- Deploy using

      gcloud app deploy --project dinodent-forms
