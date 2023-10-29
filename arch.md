@startuml
start

:Get a [Website] to process;

:Navigate to the baseUrl on the [Website];

:Run any login steps (from the [Website]);

if (Has the home page been processed recently?) then (no)
:Process the home page;
endif

: Get a set of [Page] to process;

while (Have more [Page] to process?) is (yes)
:Process each page;

    :Get all the content, title, and keywords from the page;

    :Save any new links in the [Page] table (with no-content);

    :Get another set of [Page];

endwhile

:Return the results of the processing;

stop
@enduml
