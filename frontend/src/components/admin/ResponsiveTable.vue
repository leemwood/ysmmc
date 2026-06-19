<script setup lang="ts" generic="T extends Record<string, any>">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { cn } from '@/lib/utils'

export interface ColumnConfig<T> {
  key: string
  title: string
  class?: string
  width?: string
  align?: 'left' | 'center' | 'right'
  formatter?: (row: T) => string
}

interface Props {
  columns: ColumnConfig<T>[]
  rows: T[]
  loading?: boolean
  keyField?: keyof T
  skeletonRows?: number
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  keyField: 'id' as keyof T,
  skeletonRows: 5,
})

function getRowKey(row: T): string | number {
  const value = (row as Record<string, any>)[props.keyField as string]
  return value !== undefined ? String(value) : Math.random().toString(36).slice(2)
}

function getCellValue(row: T, column: ColumnConfig<T>): string {
  if (column.formatter) {
    return column.formatter(row)
  }
  const value = row[column.key]
  return value !== undefined && value !== null ? String(value) : ''
}

function getAlignClass(align?: string): string {
  switch (align) {
    case 'center': return 'text-center'
    case 'right': return 'text-right'
    default: return 'text-left'
  }
}
</script>

<template>
  <div :class="props.class">
    <div v-if="loading" class="hidden sm:block">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead v-for="column in columns" :key="column.key" :style="column.width ? { width: column.width } : undefined">
              {{ column.title }}
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="i in skeletonRows" :key="`sk-${i}`">
            <TableCell v-for="column in columns" :key="column.key">
              <Skeleton class="h-4 w-3/4" />
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <div v-else-if="rows.length === 0" class="py-12 text-center text-muted-foreground">
      <slot name="empty">
        暂无数据
      </slot>
    </div>

    <div v-else class="hidden sm:block">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead
              v-for="column in columns"
              :key="column.key"
              :class="cn(column.class, getAlignClass(column.align))"
              :style="column.width ? { width: column.width } : undefined"
            >
              {{ column.title }}
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="row in rows" :key="getRowKey(row)">
            <TableCell
              v-for="column in columns"
              :key="column.key"
              :class="cn(column.class, getAlignClass(column.align))"
            >
              <slot :name="column.key" :row="row">
                {{ getCellValue(row, column) }}
              </slot>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <div class="sm:hidden space-y-3">
      <template v-if="loading">
        <Card v-for="i in skeletonRows" :key="`m-sk-${i}`">
          <CardContent class="p-4 space-y-3">
            <Skeleton class="h-5 w-1/2" />
            <Skeleton class="h-4 w-full" />
            <Skeleton class="h-4 w-2/3" />
          </CardContent>
        </Card>
      </template>
      <template v-else>
        <Card v-for="row in rows" :key="`m-${getRowKey(row)}`" class="overflow-hidden">
          <CardContent class="p-4">
            <slot name="mobile-card" :row="row">
              <div class="space-y-2">
                <div
                  v-for="column in columns.filter(c => c.key !== 'actions')"
                  :key="column.key"
                  class="flex justify-between gap-2"
                >
                  <span class="text-sm text-muted-foreground">{{ column.title }}</span>
                  <span class="text-sm font-medium text-right">
                    <slot :name="column.key" :row="row">
                      {{ getCellValue(row, column) }}
                    </slot>
                  </span>
                </div>
              </div>
            </slot>
            <div v-if="$slots.actions" class="mt-3 flex flex-wrap gap-2 justify-end">
              <slot name="actions" :row="row" />
            </div>
          </CardContent>
        </Card>
      </template>
    </div>
  </div>
</template>
