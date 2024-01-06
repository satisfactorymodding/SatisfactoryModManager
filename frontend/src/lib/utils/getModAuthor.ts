export function getAuthor(mod?: { authors: { role: string, user: { username: string } }[] } | null): string | undefined {
  return mod ? mod.authors.filter((author) => author.role === 'creator')[0]?.user?.username : undefined;
}