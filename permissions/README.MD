
data, read
data, write
data_i18n, write

data_sensitive, read
data_sensitive, write
data_sensitive_i18n, write

data_private, read
data_private, write
data_private_i18n, write

metadata, read
metadata, write
metadata_i18n, write

upload
tag

data
    collections
    content_units
    files
    operations

metadata
    authors
    sources
    tags
    persons
    publishers

archive_admin       all
archive_editor      data_sensitive_*, metadata*, upload, tag
archive_uploader    data_sensitive_*, upload
archive_typist      data* read, upload
archive_tagger      data_sensitive_*, tag
bb_user             data read, metadata read

Pending Approval
BB Users
    IT
    Finance
    Shidur
    Management
    Akladot
        Maklidot
        Editors
