# Requirements
 - `docker-compose`
 - GO

# ORMs
## What Wikipedia says
> Object-relational mapping (ORM, O/RM, and O/R mapping tool) in computer science is a programming technique for converting data between incompatible type systems using object-oriented programming languages. This creates, in effect, a "virtual object database" that can be used from within the programming language

## The need for an ORM
 - abstracts complexity of the underlying data storage
    * how do I connect to the database?
    * do I need to keep connections open?
    * I need to run a query which will return millions of rows, how do I do that?
    * how do I define relationships between different types of data?
 - easier to reason about
    * most programming languages have at least one well-maintained ORM
    * the APIs feel familiar
 - flexibility
    * switching the data storage layer (portable code)

## The bad
 - opinionated
 - performance vs. simplicity

## Some famous ORMs
 - Django's (python)
 - Rails' (ActiveRecord, Ruby)
 - Hibernate (Java)
 - Diesel (rust)
 - TypeORM (typescript, javascript)

# GORM
> The fantastic ORM library for Golang

## The basics
 - migrations (`./connection` and `./migrations`)
 - writing (`./basics`)
 - reading (`./basics`)
 - pagination (`./basics`)

## More meaty stuff
 - joins (left, right, inner, etc.) (`./joins`)
 - transactions (`./transactions`)
 - explicit locking (`./locks`)
 - connection pools (`./pools`)

# As an ECHO Engineer
Now you have a pretty good idea of how GORM works and when an ORM like this is useful, there are a few things you always need to remember. 

## While writing code that requires any sort of interaction with a database
  - security: make sure you are familiar with SQL injection and how to avoid it. Be very careful with just passing un-escaped user input to your queries. **ALWAYS** parametrise queries. In GORM: `conn.Where("name = ?", userInput)` and **NEVER** `conn.Where("name = "+userInput)`
 - think carefully about your data. Take your time to decide how you want to map business concepts into database tables. A schema that sounds reasonable and 'natural' in your brain might be terrible if you translate it literarlly to tables in your database. Always consider:
   * who is going to query it?
   * do I need to optimise for writes, reads or both? do I even need to care about optimisations?
   * how many rows am I expected to have in each table?
   * normalisation: how much is too much?
   * do I need to support concurrent access?
 - indexes tend to be used as a swiss army knife. While in many cases it is just fine to add one, make sure you understand the implications:
   * they are not free. Indexes need to be updated as your tables grow, which can impact write-sensitive applications
   * the query planner will always outsmart you. Make sure you understand how your queries run in case they are already using indexes when you think they are not (see below)
 - `EXPLAIN` and `EXPLAIN ANALYZE` are your friends. If you are unsure of whether a query you are writing is slow or not, run it with these directives (eg: `EXPLAIN SELECT * FROM ...`). They will show you the query plan, estimated cost (`EXPLAIN`) and actual cost (`EXPLAIN ANALYZE`). These are extremely useful to debug the performance of your queries and help you decide if you need to tweak your query, change your data model, add new indexes, etc.
 - **be very careful** with migrations:
   * as of today, we run database migrations when our services boot up
   * a migration that cannot be applied will cause downtime in our service. The best way to get out of this situation is to roll back the service to a known previous working state. Always make sure that rolling back will be safe
   * migrations are applied in order, based on the migration's filename. Make sure you name them correctly (`1_add_foo.up.sql`, `2_add_bar.up.sql`...`N_alter_foo_column.up.sql`). Remember the `up` (upward migration). Ideally we'd have a `down` counterpart (which reverts the change). We currently don't support this though. See [this](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) for more information
   * think about the size of the table where you're applying the migration. **Specially when you are adding a new field**. Does it need a default value?. Sometimes you can work around this on the application code. If you can't, make the field NULLABLE and then write a data migration or a script that populates it (from your app). If you try to add a new column with a default value in a table containing hundreds of thousands of rows (`orders`, `patient_medications`, etc.), the database will lock the target table while the column with the default value is written to all existing rows. This can cause downtime
 - in general, try to avoid explicit locking. Tables holding references to other tables are specially sensitive to this. If you need to lock some rows, make sure you understand the different types of locks you can use and how they interact between each other. Always lock as fewest as possible. In some cases, a really well thought data model along with a proper system design might avoid all sorts of explicit locking. In others, the complexity you'd have to add to the system to avoid locking would be so high that it's not worth trying to avoid them. Optimistic locking is also something to consider. It's all about the trade-offs.
 - ORMs are great, until they are not. If your query involves several complex JOINs, don't assume GORM will generate the best performing query for you. If you really care about performance, dig into what query or queries GORM is making. If you just care about funcionality, then in most cases you'll be fine with GORM

## If you have production access
 - avoid doing direct queries to them as much as you can
 - `BEGIN` is your **best** friend. This will make any changes you make reversible (`ROLLBACK`) or definitive (`COMMIT`)
 - never try new queries in a production database. If you need a big dataset in order to benchmark a query, build that dataset locally

## And of course
If you are not sure about something, **ask for help!**