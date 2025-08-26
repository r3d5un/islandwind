const options: Intl.DateTimeFormatOptions = {
  year: 'numeric',
  month: 'long',
  day: 'numeric',
}

export function formatDate(date: Date): string {
  return date.toLocaleDateString(undefined, options)
}
