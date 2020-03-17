import React from 'react';
import {
    List,
    Edit,
    Create,
    Datagrid,
    TextField,
    ReferenceField,
    EditButton,
    SimpleForm,
    SelectInput,
    ReferenceInput,
    TextInput,
    Filter
} from 'react-admin';
const PostFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Search" source="q" alwaysOn />
        <ReferenceInput label="User" source="userId" reference="users" allowEmpty>
            <SelectInput optionText="name" />
        </ReferenceInput>
    </Filter>
);
export const PostList = props => (
    <List {...props} filters={<PostFilter />}>
        <Datagrid>
            <TextField source={"id"} />
            <ReferenceField source="userId" reference="users">
                <TextField source="name" />
            </ReferenceField>
            <TextField source="title" />
            <EditButton/>
        </Datagrid>
    </List>
);
const PostTitle = ({ record }) => {
    return <span>Post {record ? `"${record.title}"` : ''}</span>;
};
export const PostEdit = props => (
    <Edit {...props} title={<PostTitle/>}>
        <SimpleForm>
            <ReferenceField source={"userId"} reference={"users"}>
                <SelectInput source={"name"} />
            </ReferenceField>
            <TextField source={"title"} />
            <TextField source={"body"} multiline />
        </SimpleForm>
    </Edit>
);

export const PostCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <ReferenceInput source="userId" reference="users">
                <SelectInput optionText="name" />
            </ReferenceInput>
            <TextInput source="title" />
            <TextInput multiline source="body" />
        </SimpleForm>
    </Create>
);