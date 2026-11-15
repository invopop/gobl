## Fixed

- `net`: `Client.VerifyDelivery`, `Client.Who`, `Client.VerifyEndorsement` and `Client.FindVerifications` reject an envelope with no document, or whose document does not match its digest.
