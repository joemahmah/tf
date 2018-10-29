# tf - A File Tagger Utility
tf is a simple file tagger.


## Usage
tf is currently a CLI utility, but may be expanded at some point to have a tui and/or gui.

CLI tf uses subcommands to determine what you want to do. `tf <subcommand> <args>`

The following subcommands are available:
* tag - Adds tags to files. You can add as many tags as you want to as many files as you want.
* untag - Removes tags from files. You can specify tags to remove or just remove them all. Untracks a file if said file has no remaining tags.
* list - Lists the tags associated with a specified file.
* ls - Lists the files in a directory. All files include tags. Untagged files are also listed.
* lstag - Lists the files associated with a specific tag. Note: future plans to add flag to list all flags registered.
* mvtag - Renames a tag (not implemented)
* cptag - Copies a tag (not implemented)
* rmtag - Removes a tag (not implemented)
* merge - Merges two tags together (not implemented)
* mv - Moves a file while keeping tags.
* cp - Copes a file and its tags.
* rm - Removes a file and untracks it.

### File Manipulation
tf uses a database to store data on files rather than metadata in the files. Because of this, moving and copying files will result in the new file not being tagged.

Likewise, deleting a file will keep the tags for that file rather than also deleting them.

If you want to move, copy, or delete files, you will want to use a tf subcommand to do so.

## Installation
Currently, there is not executable available for download. Executable downloads are planned for Windows and Linux systems.

You can manually compile and install with `go install -u github.com/joemahmah/tf`. This option will require a go compiler.
