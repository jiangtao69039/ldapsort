## LDAP SORT

#### Description

The data retrieved from the LDAP server using the ldapsearch command is unordered, and the retrieved data is not stored in DN order. Therefore, when importing back, it usually fails because there is a sequential dependency between DNs.

#### Example:
```
DN: ou=OU1

DN: people=admin001,ou=OU1
```

The above two entries must be imported in the first order when importing to LDAP.

#### Function
Sort a .ldif data file in DN order, and the generated file can be correctly imported into LDAP.

#### Usage
Download the ldapsort_linux_x86 binary from the release.
Change the file permissions to make it executable with `chmod +x ldapsort_linux_x86`.
Rename the *.ldif file to data.ldif and place it in the same directory as ldapsort_linux_x86.
Execute `./ldapsort_linux_x86` to generate a data_sort.ldif file.
The data_sort.ldif is the sorted file.
