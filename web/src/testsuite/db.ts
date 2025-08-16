import { Client } from 'pg'
import { promises as fs } from 'fs'
import { join } from 'path'
import { readFile } from 'node:fs/promises'

export class MigrationError extends Error {
  constructor(public message: string = 'Unable to run migraiton') {
    super(message)
    this.name = 'MigrationError'
  }
}

export enum MigrationDirection {
  up = 'up',
  down = 'down',
}

export async function migrate(client: Client, direction: MigrationDirection): Promise<void> {
  const migrationDirectory: string = './../migrations'
  try {
    const files = await fs.readdir(migrationDirectory)
    const upSqlFiles = files.filter((file) => file.endsWith(`.${direction}.sql`))

    for (const file of upSqlFiles) {
      const filePath = join(migrationDirectory, file)
      const statement: string = await readFile(filePath, { encoding: 'utf-8' })
      console.log(statement)
      await client.query(statement)
    }
  } catch {
    throw new MigrationError()
  }
}
