# term-apply

term-apply is a custom SSH server used by candidates to submit their resumes

## Usage

term-apply is maintained and distriubuted as a [nix package](https://github.com/Nebulaworks/nix-garage/blob/master/pkgs/term-apply/default.nix) and [container](https://hub.docker.com/r/nebulaworks/term-apply).

## DynamoDB

term-apply leverages DynamoDB for data storage and retrieval. The general schema for each applicant is as follows:

> applied_date: unix time - number - Sort DDB Key <br>
> email: candy@date.com - string - Primary/Partition DDB Key <br>
> name: Candy Date - string <br>
> github: candydate100 - string <br>
> role_applied: sr. software engineer - string <br>
> role_override: string - in the case were their role is different otherwise null <br>
> resume_review: bool - Resume passed or fail the review <br>
> resume_review_date: number - resume reviewed date <br>
> interviews: Map <br>
>   1: {date: number, pass: bool, notes: binary} <br>
>   2: {date: number, pass: bool, notes: binary} <br>
> ignore_workflow: bool - ignore from automation workflow <br>
> offer_given: bool <br>
> offer_date: number <br>
> offer_accepted:  bool <br>
> offer_accepted_date: number <br>
> rejected: bool <br>
> rejected_date: number <br>
> rejected_msg_override: binary - Message custom to applicant  <br>

## Development

These instructions recommend using `nix-shell`. If you choose not to, please make sure you have a functional `go 1.17` installation and the `make` command installed.

1. You will need access to an S3 bucket and prefix(es) with the following object permissions.

```
s3:GetObject
s3:PutObject
s3:PutObjectTagging
```

2. From the directory where this file is located, enter a nix shell
```
nix-shell
```

3. Set the `TA_BUCKET` variable
```
export TA_BUCKET=my-bucket
```

4. If necessary, set other [environment variables](#environment-variables) for your specific environment.
> For most usage, the defaults are fine for local development <br><br>
> You will also need to make sure that the default AWS configuration has read/write access to your s3 bucket and the region is set correctly. 

5. Create an `uploads` directory if one does not yet exist
```
mkdir uploads
```

6. Code

7. Compile the code
```
make test && make build
```

8. Test
```
./term-apply
```

## Environment Variables

| Name | Description | Default |
|------|-------------|---------|
| TA_BUCKET | name of the S3 bucket. (ex: `my-bucket`) This is the only **required** field and has no sane default. We opt to fail vs accidently using an incorrect bucket. | "" |
| TA_HOST | the interface IP to listen on  | "0.0.0.0" |
| TA_PORT | the TCP port to listen on | 23234 |
| TA_UPLOAD_DIR | the path where temp resumes will be stored before being sent to S3 | "./uploads" |
| TA_DATAFILE | the csv file name (both locally and in S3) | "applicants.csv" |
| TA_CSV_PREFIX | the S3 prefix where the `TA_DATAFILE` will be stored | "/term-apply/dev/data" |
| TA_RESUME_PREFIX | the S3 prefix where the uploaded PDFs will be stored | "/term-apply/dev/resumes" |
| TA_SSM_HOST_KEY_PARAM | name of the SSM Parameter that holds the ssh host key for the runtime environment. If none is given, a host key is automatically generated | "" |
| TA_HOST_KEY_PATH | the local path where the ssh host key is located and where it will be generated if no key exists at this location. If an SSM Parameter is provided, this is also the target download location for the stored key. | ".ssh/term_info_ed25519" |
