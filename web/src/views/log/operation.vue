<template>
  <div style="width: 100%; height:100%;"/>
  <TableBody>
    <template #header>
      <TableHeader ref="tableHeaderRef" :show-filter="false">
        <template #table-config :title="'操作日志'">
          <a-form :model="{}" layout="inline">
            <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
              <template v-if="item.render">
                <FormRender :formItem="item" :render="item.render"/>
              </template>
              <template v-else>
                <template v-if="item.type === 'input'">
                  <a-input v-model="item.value.value" :placeholder="item.placeholder"/>
                </template>
                <template v-if="item.type === 'select'">
                  <a-select
                      v-model="item.value.value"
                      :placeholder="item.placeholder"
                      style="width: 150px"
                  >
                    <a-option
                        v-for="optionItem of item.optionItems"
                        :key="optionItem.value"
                        :value="optionItem.value"
                    >
                      {{ optionItem.label }}
                    </a-option>
                  </a-select>
                </template>
                <template v-if="item.type === 'date'">
                  <a-date-picker v-model="item.value.value"/>
                </template>
                <template v-if="item.type === 'time'">
                  <a-time-picker v-model="item.value.value" value-format="HH:mm:ss"/>
                </template>
                <template v-if="item.type === 'check-group'">
                  <a-checkbox-group v-model="item.value.value">
                    <a-checkbox v-for="it of item.optionItems" :key="it.value" :value="it.value">
                      {{ item.label }}
                    </a-checkbox>
                  </a-checkbox-group>
                </template>
              </template>
            </a-form-item>
          </a-form>
          <a-button size="small" type="primary" @click="onSearch">查询</a-button>
        </template>
      </TableHeader>
    </template>
    <template #default>
      <a-table
          :column-resizable="false"
          :data="dataList"
          :hoverable="true"
          :loading="tableLoading"
          :pagination="false"
          :rowKey="rowKey"
          :scroll="{ y: tableHeight }"
          :stripe="true"
          show-header
          size="small"
          table-layout-fixe
          @selection-change="onSelectionChange"
      >
        <template #columns>
          <a-table-column
              v-for="item of tableColumns"
              :key="item.identity"
              :align="item.align"
              :data-index="(item.key as string)"
              :fixed="item.fixed"
              :title="(item.title as string)"
              :width="item.width"
          >
            <template v-if="item.key === 'index'" #cell="{ rowIndex }">
              {{ rowIndex + 1 }}
            </template>
            <template v-else-if="item.key === 'http_code'" #cell="{ record }">
              <a-tag :color="record.http_code ==200 ? 'green' : 'red'">
                {{ record.http_code === 200 ? '成功' : '失败' }}
              </a-tag>
            </template>
            <template v-else-if="item.key === 'time'" #cell="{ record }">
              <a-tag :color="record.time <=100 ? 'blue' : 'red'">
                {{ record.time + 'ms' }}
              </a-tag>
            </template>
            <template v-else-if="item.key === 'method'" #cell="{ record }">
              <a-tag :color="record.method ==='GET' ? 'purple' : 'pink'">
                {{ record.method }}
              </a-tag>
            </template>
            <template v-else-if="item.key === 'image_Url'" #cell="{}">
              <a-avatar :size="30" :style="{ backgroundColor: 'var(--color-primary-light-1)' }">
                <IconUser/>
              </a-avatar>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </template>

    <template #footer>
      <TableFooter ref="tableFooterRef" :pagination="pagination"/>
    </template>
  </TableBody>
  <!--    继续-->
</template>

<script lang="ts">
import {logList,} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn, useTableHeight,} from '@/hooks/table'
import {Form, Input, Message} from '@arco-design/web-vue'
import {defineComponent, getCurrentInstance, h, onMounted, ref} from 'vue'
import AddButton from "@/components/AddButton.vue";
import useUserStore from "@/store/modules/user";
import {FormItem, ModalDialogType} from "@/types/components";
import usePost from '@/hooks/usePost'
import useGet from "@/hooks/useGet";
import FormRender from "@/components/FormRender";

export default defineComponent({
  components: {FormRender, AddButton},
  name: 'UserList',
  setup: function () {
    const table = useTable()
    const rowKey = useRowKey('id')
    const pagination = usePagination(doRefresh)
    const {selectedRowKeys, onSelectionChange} = useRowSelection()


    const conditionItems: Array<FormItem> = [
      {
        key: 'IP',
        label: 'IP地址',
        type: 'input',
        placeholder: '请输入IP地址',
        value: ref(''),
        reset: function () {
          this.value.value = ''
        },
        render: (formItem: FormItem) => {
          return h(Input, {
            placeholder: '请输入IP地址',
            modelValue: formItem.value.value,
            'onUpdate:modelValue': (value) => {
              formItem.value.value = value
            },
          })
        },
      },
      {
        key: 'method',
        label: '请求方式',
        value: ref(),
        type: 'select',
        placeholder: '请选择请求方式',
        optionItems: [
          {
            label: 'GET',
            value: 'GET',
          },
          {
            label: 'POST',
            value: 'POST',
          },
        ],
        reset: function () {
          this.value.value = undefined
        },
      },
    ]


    const formRef = ref<typeof Form>()


    const formModel = ref({})
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const actionTitle = ref('添加员工')
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '操作者',
        key: 'uid',
        dataIndex: 'uid',
        width: 100,
      },
      {
        title: '请求方式',
        key: 'method',
        dataIndex: 'method',
      },
      {
        title: '请求URL',
        key: 'path',
        dataIndex: 'path',
        width: 300,
        align: 'left'
      },
      {
        title: 'IP地址',
        key: 'ip',
        dataIndex: 'ip ',
        align: 'left',
        width: 350,
      },
      {
        title: '耗时',
        key: 'time',
        dataIndex: 'time',
        width: 180,
      },
      {
        title: '状态',
        key: 'http_code',
        dataIndex: 'http_code',
        width: 100,

      },
      {
        title: '操作时间',
        key: 'CreatedAt',
        dataIndex: 'CreatedAt',
      },
    ])
    const expandAllFlag = ref(true)
    const userStore = useUserStore()
    const post = usePost()
    const get = useGet()


    function doRefresh() {
      get({
        url: logList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            offset: pagination.page,
            limit: pagination.pageSize,
          }
        },
      }).then((res) => {
        res.data.forEach(i => {
          const date = new Date(i.CreatedAt);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          i.CreatedAt = `${year}-${month}-${day}  ${date.getHours() >= 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() >= 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() >= 10 ? date.getSeconds() : '0' + date.getSeconds()}`
        })
        table.handleSuccess(res)
        console.log(res)
        pagination.setTotalSize(res.count)
      }).catch(error => {
        Message.error(error.message)
        table.tableLoading.value = false
      })
    }

    function onSearch() {
      let data: any = conditionItems.reduce((pre, cur) => {
        ;(pre as any)[cur.key] = cur.value.value
        return pre
      }, {})
      const tableList = table.dataList.filter(i => {
        if (i.method == data.method) {
          return i
        }
      })
      table.handleSuccess({data: tableList})
      pagination.setTotalSize(tableList.length || 10)
    }


    onMounted(async () => {
      table.tableHeight.value = await useTableHeight(getCurrentInstance())
      doRefresh()
    })
    return {
      ...table,
      rowKey,
      selectedRowKeys,
      onSelectionChange,
      onSearch,
      expandAllFlag,
      tableColumns,
      pagination,
      formModel,
      actionTitle,
      modalDialogRef,
      formRef,
      conditionItems
    }
  },
})
</script>
